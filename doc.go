/*
Package realip provides cliant real ip detection mechanisms from http.Request
in Go's HTTP middleware layer.

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
	l, _ := net.Listen("tcp", "127.0.0.1:8080")

	http.Serve(l, handler)

This realizes the same function as Nginx's ngx_http_realip_module in the layer inside Go.
Therefore, the setting property names and behaviors are also close to ngx_http_realip_module.

The realip provides Go's HTTP Middleware. It detects the client's real ip from the request and
sets it in the specified request header (Default: X-Real-IP).
*/
package realip
