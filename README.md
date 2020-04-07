go-realip
=======

[![Test Status](https://github.com/natureglobal/go-realip/workflows/test/badge.svg?branch=master)][actions]
[![Coverage Status](https://coveralls.io/repos/natureglobal/go-realip/badge.svg?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/natureglobal/go-realip?status.svg)][godoc]

[actions]: https://github.com/natureglobal/go-realip/actions?workflow=test
[coveralls]: https://coveralls.io/r/natureglobal/go-realip?branch=master
[license]: https://github.com/natureglobal/go-realip/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/natureglobal/go-realip

The go-realip detects RealIP in Go's http middleware layer.

## Synopsis

```go
_, ipnet, _ := net.ParseCIDR("192.168.0.0/16")
var middleware func(http.Handler) http.Handler = realip.MustMiddleware(&realip.Config{
    RealIPFrom:      []*net.IPNet{ipnet},
    RealIPHeader:    realip.HeaderXForwardedFor,
    RealIPRecursive: true,
})
var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World")
})
handler = middleware(handler)
```

## Description

The go-realip package implements detect clietn real IP from request headers in Go's http middleware layer.
This have a similar behavior as Nginx's ngx\_http\_go-realip\_module

## Installation

```console
% go get github.com/natureglobal/go-realip
```

## Author

- [Songmu](https://github.com/Songmu)
