package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	_ "net"
	"net/http"
	_ "net/http"
	"net/netip"
	"reflect"
	"strings"
	"unsafe"
)

//go:linkname reflect_typelinks reflect.typelinks
func reflect_typelinks() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname reflect_add reflect.add
func reflect_add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer

func main() {
	var url string
	var ip string

	flag.StringVar(&url, "url", "http://www.google.com", "url to send query")
	flag.StringVar(&ip, "ip", "", "override ip address")
	flag.Parse()

	ctx := context.TODO()
	var useGo bool
	if ip != "" {
		addr, err := netip.ParseAddr(ip)
		if err != nil {
			log.Panicln("failed to parse ip", ip, err)
		}

		ipaddr := net.IPAddr{
			IP:   addr.AsSlice(),
			Zone: addr.Zone(),
		}
		var keytype reflect.Type
		sections, offsets := reflect_typelinks()
	bmain:
		for i, base := range sections {
			for _, offset := range offsets[i] {
				typeAddr := reflect_add(base, uintptr(offset), "")
				typ := reflect.TypeOf(*(*interface{})(unsafe.Pointer(&typeAddr)))

				if strings.Contains(typ.String(), "nettrace.LookupIPAltResolverKey") {
					keytype = typ
					break bmain
				}
			}

		}
		if keytype == nil {
			log.Panicln("failed to find nettrace.LookupIPAltResolverKey context key, unable to inject custom resolver")
		}
		key := reflect.Indirect(reflect.New(keytype.Elem())).Interface()
		ctx = context.WithValue(ctx, key, func(ctx context.Context, network string, name string) ([]net.IPAddr, error) {
			return []net.IPAddr{ipaddr}, nil
		})
		useGo = true
	}

	dialer := net.Dialer{Resolver: &net.Resolver{
		PreferGo: useGo,
	}}
	t := http.DefaultTransport.(*http.Transport)
	t.DialContext = dialer.DialContext

	req, _ := http.NewRequest(http.MethodHead, url, nil)
	req = req.WithContext(ctx)
	res, err := t.RoundTrip(req)

	if err != nil {
		log.Panicln("failed to send head request", url, err)
	}

	fmt.Println(res.Status, res.Header)
}
