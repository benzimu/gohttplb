GoHTTPLB
========

GoHTTPLB: load balancing http client for Go

## Load balancing strategy

There are three strategies for GoHTTPLB now, default `PolicyRandom`:

* `PolicyRandom` - random selection service instance. If return err, the remaining service instances will be retried randomly.
* `PolicyOrder` - order selection service instance. If rerurn err, the remaining service instances will be retried in order.
* `PolicyFlow` - select service instance by flow. `TODO`

## Installation

```bash
$ go get github.com/beeeeeeenny/gohttplb
```

## Usage

```go
lbclient := gohttplb.New()
```
