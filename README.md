realip
=======

[![Test Status](https://github.com/natureglobal/realip/workflows/test/badge.svg?branch=master)][actions]
[![Coverage Status](https://coveralls.io/repos/natureglobal/realip/badge.svg?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/natureglobal/realip?status.svg)][godoc]

[actions]: https://github.com/natureglobal/realip/actions?workflow=test
[coveralls]: https://coveralls.io/r/natureglobal/realip?branch=master
[license]: https://github.com/natureglobal/realip/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/natureglobal/realip

The realip detects client real ip in Go's HTTP middleware layer.

## Synopsis

```go
_, ipnet, _ := net.ParseCIDR("192.168.0.0/16")
var middleware func(http.Handler) http.Handler = realip.MustMiddleware(&realip.Config{
    RealIPFrom:      []*net.IPNet{ipnet},
    RealIPHeader:    realip.HeaderXForwardedFor,
    RealIPRecursive: true,
})
var handler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, req.Header.Get("X-Real-IP"))
})
handler = middleware(handler)
```

## Description

The realip package implements detecting client real IP mechanisms from request headers in Go's HTTP middleware layer.

This realizes the similar function as Nginx's `ngx_http_realip_module` in the layer inside Go.
Therefore, the setting property names and behaviors are also close to `ngx_http_realip_module`.

The realip provides Go's HTTP Middleware. It detects the client's real ip from the request and
sets it in the specified request header (Default: X-Real-IP).

## Installation

```console
% go get github.com/natureglobal/realip
```

## Author

- [Songmu](https://github.com/Songmu)
