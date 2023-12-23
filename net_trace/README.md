# Net Trace

Injecting Golang Net Tracer into dial context.

## Running

```golang
$ go run . -url http://google.com
2023/12/23 16:34:45 Trace: DNSStart google.com
2023/12/23 16:34:45 Trace: DNSDone [{2404:6800:4001:801::200e } {142.251.222.238 }] false <nil>
2023/12/23 16:34:45 Trace: ConnectStart tcp [2404:6800:4001:801::200e]:80
2023/12/23 16:34:45 Trace: ConnectDone tcp [2404:6800:4001:801::200e]:80 <nil>
301 Moved Permanently map[Cache-Control:[public, max-age=2592000] ...
```
