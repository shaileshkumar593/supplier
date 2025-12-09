package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// NewIPv4Client enforces IPv4 and forces HTTP/1.1 end-to-end with hardened timeouts/pooling.
// NewIPv4Client enforces IPv4 and logs resolution/dial/TLS when IPV4_DEBUG=1.
// Signature unchanged.
func NewIPv4Client(timeout time.Duration, proxy *url.URL, ipv4DirectDomains []string) *http.Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	// opt-in debug logger (stdout) without changing signature
	type dbg interface{ Printf(string, ...any) }
	var d dbg
	if os.Getenv("IPV4_DEBUG") == "1" {
		d = log.New(os.Stdout, "ipv4client ", log.LstdFlags|log.Lmicroseconds)
	} else {
		d = log.New(io.Discard, "", 0)
	}

	shouldBypass := func(host string) bool {
		h := strings.ToLower(host)
		for _, s := range ipv4DirectDomains {
			if strings.HasSuffix(h, strings.ToLower(s)) {
				return true
			}
		}
		return false
	}

	resolver := net.DefaultResolver
	baseDialer := &net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}

	resolveIPv4 := func(ctx context.Context, host string) ([]net.IP, error) {
		start := time.Now()
		ips, err := resolver.LookupIP(ctx, "ip4", host)
		if err != nil {
			d.Printf("dns ip4 error host=%q err=%v", host, err)
			return nil, err
		}
		d.Printf("dns ip4 done host=%q ips=%v dur=%s", host, ips, time.Since(start))
		if len(ips) == 0 {
			return nil, fmt.Errorf("no IPv4 A record for %s", host)
		}
		return ips, nil
	}

	dialIPv4 := func(ctx context.Context, host, port string) (net.Conn, error) {
		ips, err := resolveIPv4(ctx, host)
		if err != nil {
			return nil, err
		}
		var last error
		for _, ip := range ips {
			addr := net.JoinHostPort(ip.String(), port)
			d.Printf("dial start addr=%s", addr)
			c, err := baseDialer.DialContext(ctx, "tcp4", addr)
			if err == nil {
				d.Printf("dial ok   addr=%s local=%v", addr, c.LocalAddr())
				return c, nil
			}
			d.Printf("dial fail addr=%s err=%v", addr, err)
			last = err
		}
		return nil, last
	}

	proxyFunc := func(req *http.Request) (*url.URL, error) {
		host := req.URL.Hostname()
		if shouldBypass(host) {
			d.Printf("proxy bypass host=%q", host)
			return nil, nil
		}
		if proxy != nil {
			d.Printf("proxy custom host=%q proxy=%q", host, proxy)
			return proxy, nil
		}
		u, err := http.ProxyFromEnvironment(req)
		if err != nil {
			d.Printf("proxy env error host=%q err=%v", host, err)
			return nil, err
		}
		if u != nil {
			d.Printf("proxy env host=%q proxy=%q", host, u)
		} else {
			d.Printf("proxy none host=%q", host)
		}
		return u, nil
	}

	tr := &http.Transport{
		Proxy: proxyFunc,
		DialContext: func(ctx context.Context, _ string, address string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				d.Printf("split hostport error address=%q err=%v", address, err)
				return nil, err
			}
			return dialIPv4(ctx, host, port)
		},
		// Hardening
		TLSHandshakeTimeout:    8 * time.Second,
		ResponseHeaderTimeout:  8 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		IdleConnTimeout:        90 * time.Second,
		MaxIdleConns:           1024,
		MaxIdleConnsPerHost:    256,
		MaxConnsPerHost:        512,
		MaxResponseHeaderBytes: 1 << 20, // 1 MiB

		ForceAttemptHTTP2: true,
		TLSNextProto:      map[string]func(string, *tls.Conn) http.RoundTripper{},
		TLSClientConfig:   &tls.Config{MinVersion: tls.VersionTLS13, NextProtos: []string{"http/1.1"}},
	}

	tr.DialTLSContext = func(ctx context.Context, _ string, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			d.Printf("split host-port tls error addr=%q err=%v", addr, err)
			return nil, err
		}
		tcp, err := dialIPv4(ctx, host, port)
		if err != nil {
			return nil, err
		}
		cfg := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS13, NextProtos: []string{"http/1.1"}}
		d.Printf("tls handshake start server=%q", host)
		tlsConn := tls.Client(tcp, cfg)
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			_ = tcp.Close()
			d.Printf("tls handshake fail server=%q err=%v", host, err)
			return nil, err
		}
		cs := tlsConn.ConnectionState()
		d.Printf("tls handshake ok server=%q proto=%s vers=%x", host, cs.NegotiatedProtocol, cs.Version)
		return tlsConn, nil
	}

	return &http.Client{Transport: tr, Timeout: timeout}
}
