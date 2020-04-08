package realip_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/natureglobal/realip"
)

func mustParseCIDR(addr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(addr)
	if err != nil {
		panic(err)
	}
	return ipnet
}

func TestMiddleware(t *testing.T) {
	testCases := []struct {
		name    string
		config  *realip.Config
		headers map[string]string
		expect  string
	}{{
		name:   "X-Real-IP: default",
		expect: "127.0.0.1",
	}, {
		name: "X-Real-IP: no RealIPFrom",
		headers: map[string]string{
			"x-real-ip": "1.1.1.1",
		},
		expect: "1.1.1.1",
	}, {
		name: "X-Real-IP: not a trusted ip from",
		config: &realip.Config{
			RealIPFrom: []*net.IPNet{mustParseCIDR("192.168.0.0/16")},
		},
		headers: map[string]string{
			"x-real-ip": "1.1.1.1",
		},
		expect: "127.0.0.1",
	}, {
		name: "x-forwarded-for",
		config: &realip.Config{
			RealIPHeader: realip.HeaderXForwardedFor,
		},
		headers: map[string]string{
			"X-Forwarded-For": "192.168.0.1",
		},
		expect: "192.168.0.1",
	}, {
		name: "x-forwarded-for: recent non-trusted one",
		config: &realip.Config{
			RealIPFrom: []*net.IPNet{
				mustParseCIDR("127.0.0.1/32"),
				mustParseCIDR("192.168.0.0/16"),
			},
			RealIPHeader:    realip.HeaderXForwardedFor,
			RealIPRecursive: true,
		},
		headers: map[string]string{
			"X-Forwarded-For": "1.2.3.4, 1.1.1.1, 192.168.0.1",
		},
		expect: "1.1.1.1",
	}, {
		name: "x-forwarded-for: recent trusted one",
		config: &realip.Config{
			RealIPFrom: []*net.IPNet{
				mustParseCIDR("127.0.0.1/32"),
				mustParseCIDR("192.168.0.0/16"),
			},
			RealIPHeader: realip.HeaderXForwardedFor,
		},
		headers: map[string]string{
			"X-Forwarded-For": "1.2.3.4, 1.1.1.1, 192.168.0.1",
		},
		expect: "192.168.0.1",
	}, {
		name: "x-forwarded-for: remoteAddr is not a trusted address",
		config: &realip.Config{
			RealIPFrom: []*net.IPNet{
				mustParseCIDR("192.168.0.0/16"),
			},
			RealIPHeader: realip.HeaderXForwardedFor,
		},
		headers: map[string]string{
			"X-Forwarded-For": "1.2.3.4, 1.1.1.1, 192.168.0.1",
		},
		expect: "127.0.0.1",
	}, {
		name: "x-forwarded-for: all entries in xff is trusted ip",
		config: &realip.Config{
			RealIPFrom: []*net.IPNet{
				mustParseCIDR("127.0.0.1/32"),
				mustParseCIDR("192.168.0.0/16"),
			},
			RealIPHeader:    realip.HeaderXForwardedFor,
			RealIPRecursive: true,
		},
		headers: map[string]string{
			"X-Forwarded-For": "192.168.0.2, 192.168.0.1",
		},
		expect: "192.168.0.2",
	}, {
		name: "x-forwarded-for: no RealIPFrom config and true RealIPRecursive return left entry in xff",
		config: &realip.Config{
			RealIPHeader:    realip.HeaderXForwardedFor,
			RealIPRecursive: true,
		},
		headers: map[string]string{
			"X-Forwarded-For": "192.168.0.2, 192.168.0.1",
		},
		expect: "192.168.0.2",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := realip.MustMiddleware(tc.config)
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprintf(w, req.Header.Get(realip.HeaderXRealIP))
			}))
			ts := httptest.NewServer(handler)
			defer ts.Close()

			req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			if tc.headers != nil {
				for k, v := range tc.headers {
					req.Header.Add(k, v)
				}
			}
			r, _ := http.DefaultClient.Do(req)
			defer r.Body.Close()

			data, _ := ioutil.ReadAll(r.Body)
			out := strings.TrimSpace(string(data))
			if out != tc.expect {
				t.Errorf("out: %s, expect: %s", out, tc.expect)
			}
		})
	}

	t.Run("error", func(t *testing.T) {
		_, err := realip.Middleware(&realip.Config{
			RealIPHeader: "Forwarded",
		})
		if err == nil {
			t.Errorf("error should be occurred but nil")
		}
	})
}
