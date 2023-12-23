# Time Now

A POC to inject custom resolver using the `nettrace.LookupIPAltResolverKey` context key.
`nettrace.LookupIPAltResolverKey` is used interanally by Golang for testing.

## Running

```golang
# send HEAD request to google.com
$ go run main.go 

# send HEAD request to google.com with static IP
$ go run main.go -ip 142.251.222.238

# send HEAD request to yahoo.com with static IP
$ go run main.go -url myserver.com -ip 142.251.222.238
```
