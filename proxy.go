package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/http/httpproxy"
)

// from net/http/transport.go
var (
	// proxyConfigOnce guards proxyConfig
	envProxyOnce      sync.Once
	envProxyFuncValue func(*url.URL) (*url.URL, error)
)

// defaultProxyConfig returns a ProxyConfig value looked up
// from the environment. This mitigates expensive lookups
// on some platforms (e.g. Windows).
func envProxyFunc() func(*url.URL) (*url.URL, error) {
	envProxyOnce.Do(func() {
		envProxyFuncValue = httpproxy.FromEnvironment().ProxyFunc()
	})
	return envProxyFuncValue
}

func (s *shell) tlsDial(network, address string, cfg *tls.Config) (*tls.Conn, error) {
	u, err := envProxyFunc()(&url.URL{Scheme: "https", Host: address})
	if err == nil && u != nil {
		// u is a proxy, need to connect based on scheme
		switch u.Scheme {
		case "", "http":
			// use CONNECT - see https://en.wikipedia.org/wiki/HTTP_tunnel
			s, err := net.Dial("tcp", u.Host)
			if err != nil {
				return nil, err
			}
			// send request
			// TODO support Proxy-Authorization
			_, err = fmt.Fprintf(s, "CONNECT %s HTTP/1.1\r\n\r\n", address)
			if err != nil {
				return nil, err
			}
			// read response (we read one byte at a time, that's not great, but using bufio/etc makes a lot of things a lot more complex)
			var resp []byte
			c := make([]byte, 1)
			for {
				n, err := s.Read(c)
				if err != nil {
					s.Close()
					return nil, err
				}
				if n == 0 {
					// shouldn't happen with only 1 byte
					continue
				}
				resp = append(resp, c[0])
				if c[0] == '\n' {
					// check if resp ends in \r\n\r\n
					if bytes.HasSuffix(resp, []byte{'\r', '\n', '\r', '\n'}) {
						// we're good.
						break
					}
				}
			}
			respA := bytes.Split(resp, []byte{'\n'})
			statusline := strings.TrimSpace(string(respA[0])) // HTTP/1.1 200 OK
			statusinfo := strings.Split(statusline, " ")
			if len(statusinfo) < 2 || (statusinfo[0] != "HTTP/1.0" && statusinfo[0] != "HTTP/1.1") || statusinfo[1][0] != '2' {
				s.Close()
				return nil, fmt.Errorf("invalid status line %s", statusline)
			}

			// do TLS
			tlsconn := tls.Client(s, cfg)
			// perform handshake now
			err = tlsconn.Handshake()
			if err != nil {
				s.Close()
				return nil, err
			}
			return tlsconn, nil
		default:
			// TODO support socks4 socks5 etc?
			log.Printf("notice: unsupported proxy protcol %s, using direct connection", u.Scheme)
		}
	}
	return tls.Dial(network, address, cfg)
}
