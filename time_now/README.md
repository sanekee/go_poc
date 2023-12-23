# Time Now

A POC of overriding golang `time.Now()` function in testing using `//go:linkname` compiler directive.

## Running

```golang
$ go run .
Current Time 2023-12-23 03:23:48.524779 +0000 UTC
Is Valid? false
```

## Testing

```golang
$ go test -v ./...
=== RUN   TestIsValid
=== RUN   TestIsValid/is_valid_before_2021-01-01
=== RUN   TestIsValid/is_not_valid_after_2021-01-01
--- PASS: TestIsValid (0.00s)
    --- PASS: TestIsValid/is_valid_before_2021-01-01 (0.00s)
    --- PASS: TestIsValid/is_not_valid_after_2021-01-01 (0.00s)
PASS
```
