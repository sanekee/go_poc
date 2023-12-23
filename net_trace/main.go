package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	_ "net"
	"net/http"
	_ "net/http"
	"reflect"
	"unsafe"
)

//go:linkname reflect_typelinks reflect.typelinks
func reflect_typelinks() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname reflect_add reflect.add
func reflect_add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer

func main() {
	var url string

	flag.StringVar(&url, "url", "http://www.google.com", "url to send query")
	flag.Parse()

	ctx := newTracerContext(context.TODO())

	req, _ := http.NewRequest(http.MethodHead, url, nil)
	req = req.WithContext(ctx)
	res, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		log.Panicln("failed to send head request", url, err)
	}

	fmt.Println(res.Status, res.Header)
}

// DNSStart is called with the hostname of a DNS lookup
// before it begins.
func dnsStart(name string) {
	log.Println("Trace: DNSStart", name)
}

func newTracerContext(ctx context.Context) context.Context {
	var traceKeytype reflect.Type
	var traceType reflect.Type
	sections, offsets := reflect_typelinks()

	for i, base := range sections {
		for _, offset := range offsets[i] {
			typeAddr := reflect_add(base, uintptr(offset), "")
			typ := reflect.TypeOf(*(*interface{})(unsafe.Pointer(&typeAddr)))

			if typ.String() == "*nettrace.TraceKey" {
				traceKeytype = typ

			}
			if typ.String() == "*nettrace.Trace" {
				traceType = typ
			}
		}

	}

	if traceKeytype == nil {
		log.Panicln("failed to find nettrace.TraceKey context key")
	}

	if traceType == nil {
		log.Panicln("failed to find nettrace.Trace")
	}

	ctxKey := reflect.Indirect(reflect.New(traceKeytype.Elem())).Interface()

	tracerp := reflect.New(traceType.Elem())
	tracer := tracerp.Elem()

	for i := 0; i < tracer.Type().NumField(); i++ {
		field := tracer.Field(i)
		if field.Kind() == reflect.Func {
			if isFunc(field, reflect.ValueOf(dnsStart)) {
				field.Set(reflect.ValueOf(dnsStart))
				continue
			}
			if isFunc(field, reflect.ValueOf(dnsDone)) {
				field.Set(reflect.ValueOf(dnsDone))
				continue
			}
			if isFunc(field, reflect.ValueOf(connectStart)) {
				field.Set(reflect.ValueOf(connectStart))
				continue
			}
			if isFunc(field, reflect.ValueOf(connectDone)) {
				field.Set(reflect.ValueOf(connectDone))
				continue
			}
		}
	}

	return context.WithValue(ctx, ctxKey, tracerp.Interface())
}

// DNSDone is called after a DNS lookup completes (or fails).
// The coalesced parameter is whether singleflight de-duped
// the call. The addrs are of type net.IPAddr but can't
// actually be for circular dependency reasons.
func dnsDone(netIPs []any, coalesced bool, err error) {
	log.Println("Trace: DNSDone", netIPs, coalesced, err)
}

// ConnectStart is called before a Dial, excluding Dials made
// during DNS lookups. In the case of DualStack (Happy Eyeballs)
// dialing, this may be called multiple times, from multiple
// goroutines.
func connectStart(network, addr string) {
	log.Println("Trace: ConnectStart", network, addr)
}

// ConnectDone is called after a Dial with the results, excluding
// Dials made during DNS lookups. It may also be called multiple
// times, like ConnectStart.
func connectDone(network, addr string, err error) {
	log.Println("Trace: ConnectDone", network, addr, err)
}

func isFunc(f1 reflect.Value, f2 reflect.Value) bool {
	if f1.Kind() != reflect.Func || f2.Kind() != reflect.Func {
		return false
	}

	if f1.Type().NumIn() != f2.Type().NumIn() {
		return false
	}

	if f1.Type().NumOut() != f2.Type().NumOut() {
		return false
	}

	for i := 0; i < f1.Type().NumIn(); i++ {
		arg1 := f1.Type().In(i)
		arg2 := f2.Type().In(i)

		if arg1 != arg2 {
			return false
		}
	}

	for i := 0; i < f1.Type().NumOut(); i++ {
		arg1 := f1.Type().Out(i)
		arg2 := f2.Type().Out(i)

		if arg1 != arg2 {
			return false
		}
	}
	return true
}
