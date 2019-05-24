GoHTTPLB
========

GoHTTPLB: load balancing http client for Go

## Load balancing strategy

There are three strategies for GoHTTPLB now, default `PolicyRandom`:

* `PolicyRandom` - random selection service instance. If return err, the remaining service instances will be retried randomly.
* `PolicyOrder` - order selection service instance. If rerurn err, the remaining service instances will be retried in order.
* `PolicyFlowAverage` - select service instance by flow average policy.

## Installation

```bash
$ go get github.com/beeeeeeenny/gohttplb
```

## Usage

### Most convenient way
```go
addr := "127.0.0.1:8080,127.0.0.1:8081,127.0.0.1:8082"
lbclient, err := gohttplb.New(addr)
if err != nil {
    log.Println(err)
    return
}
```

### Use `LBConfig` init `LBClient`
```go
addr := "127.0.0.1:8080;127.0.0.1:8081;127.0.0.1:8082"
lbconf := &gohttplb.LBConfig{
    SchedPolicy: gohttplb.PolicyFlowAverage,
    Retry:       3,
    Separator:   ";",
}
lbclient, err := gohttplb.New(addr, lbconf)
if err != nil {
    log.Println(err)
    return
}
```

### Do `Get` request
```go
resp, err := lbclient.Get("/hello")
if err != nil {
    log.Println(err)
    return
}
```

### Do Json `Get` request
```go
resp, err := lbclient.JSONGet("/hello")
if err != nil {
    log.Println(err)
    return err
}
```

### Do auto parse response request

you can implement `ResponseParser` interface.
And set `LBConfig.ResponseParser` when init `LBClient`, like:
```go
addr := "127.0.0.1:8080;127.0.0.1:8081;127.0.0.1:8082"
lbconf := &gohttplb.LBConfig{
    ResponseParser: &CustomResponseParser{},
}
lbclient, err := gohttplb.New(addr, lbconf)
if err != nil {
    log.Println(err)
    return nil, err
}
```

then you can use `ParseGet`, `ParsePost`......, like:
```go
statusCode, data, err := lbclient.ParseGet("/hello")
if err != nil {
    log.Println(err)
    return
}
```

and use `JSONParseGet`, `JSONParsePost`......
```go
statusCode, data, err := lbclient.JSONParseGet("/hello")
if err != nil {
    log.Println(err)
    return
}
```
