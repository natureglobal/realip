package realip

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// Popular Headers
const (
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-IP"
	// we can also specify True-Client-IP, CF-Connecting-IP and so on to c.RealIPHeader
	headerForwarded = "Forwarded"
)

// Middleware returns a http middleware detecting the real IP of the client from request
// and set it in the request header.
func Middleware(c *Config) (func(http.Handler) http.Handler, error) {
	if c == nil {
		c = &Config{}
	}
	if strings.ToLower(c.realIPHeader()) == strings.ToLower(headerForwarded) {
		return nil, errors.New("haven't supported Forwarded header yet")
	}
	return func(next http.Handler) http.Handler {
		return c.handler(next)
	}, nil
}

// MustMiddleware returns a http middleware of realip. It panics when error occurred.
func MustMiddleware(c *Config) func(http.Handler) http.Handler {
	m, err := Middleware(c)
	if err != nil {
		panic(err)
	}
	return m
}

// Config is configuration for realip middleware.
// The fields naming of it is similar to ngx_http_realip_module.
type Config struct {
	RealIPFrom      []*net.IPNet
	RealIPHeader    string
	RealIPRecursive bool

	SetHeader string
}

func remoteIP(remoteAddr string) string {
	ip, _, _ := net.SplitHostPort(remoteAddr)
	return ip
}

func (c *Config) handler(next http.Handler) http.Handler {
	switch c.realIPHeader() {
	case HeaderXForwardedFor:
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if !trustedIP(net.ParseIP(remoteIP(req.RemoteAddr)), c.RealIPFrom) {
				if realIP := remoteIP(req.RemoteAddr); realIP != "" {
					req.Header.Set(c.setHeader(), realIP)
				}
			} else {
				realIP := realIPFromXFF(
					req.Header.Get(HeaderXForwardedFor),
					c.RealIPFrom,
					c.RealIPRecursive)
				if realIP == "" {
					realIP = remoteIP(req.RemoteAddr)
				}
				if realIP != "" {
					req.Header.Set(c.setHeader(), realIP)
				}
			}
			next.ServeHTTP(w, req)
		})
	default:
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if !trustedIP(net.ParseIP(remoteIP(req.RemoteAddr)), c.RealIPFrom) {
				if realIP := remoteIP(req.RemoteAddr); realIP != "" {
					req.Header.Set(c.setHeader(), realIP)
				}
			} else {
				realIP := req.Header.Get(c.realIPHeader())
				if realIP == "" {
					realIP = remoteIP(req.RemoteAddr)
				}
				if realIP != "" {
					req.Header.Set(c.setHeader(), realIP)
				}
			}
			next.ServeHTTP(w, req)
		})
	}
}

func trustedIP(ip net.IP, realIPFrom []*net.IPNet) bool {
	if len(realIPFrom) == 0 {
		return true
	}
	for _, fromIP := range realIPFrom {
		if fromIP.Contains(ip) {
			return true
		}
	}
	return false
}

func (c *Config) realIPHeader() string {
	if c.RealIPHeader == "" {
		return HeaderXRealIP // default
	}
	return c.RealIPHeader
}

func (c *Config) setHeader() string {
	if c.SetHeader == "" {
		return HeaderXRealIP // default
	}
	return c.SetHeader
}

func realIPFromXFF(xff string, realIPFrom []*net.IPNet, recursive bool) string {
	ips := strings.Split(xff, ",")
	if len(ips) == 0 {
		return ""
	}
	if !recursive {
		return strings.TrimSpace(ips[len(ips)-1])
	}
	for i := len(ips) - 1; i >= 0; i-- {
		ipStr := strings.TrimSpace(ips[i])
		if !trustedIP(net.ParseIP(ipStr), realIPFrom) {
			return ipStr
		}
	}
	return strings.TrimSpace(ips[0])
}
