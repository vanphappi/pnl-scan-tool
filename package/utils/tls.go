package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 30 * time.Second, // Set a timeout for the entire request
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true, // Prioritize server's cipher suite order
			CurvePreferences: []tls.CurveID{
				tls.CurveP256, tls.X25519, // Strong elliptic curves
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		},
		MaxIdleConns:        100,              // Pool idle connections
		MaxIdleConnsPerHost: 10,               // Per-host connection limit
		IdleConnTimeout:     90 * time.Second, // Timeout for idle connections
		MaxConnsPerHost:     20,               // Limit the maximum simultaneous connections to a host
		DisableKeepAlives:   false,            // Enable keep-alive for better performance
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Timeout for dialing the connection
			KeepAlive: 30 * time.Second, // Keep-alive settings
		}).DialContext,
		ExpectContinueTimeout: 1 * time.Second,           // Optimize for HTTP/1.1 Continue handling
		TLSHandshakeTimeout:   10 * time.Second,          // Set timeout for TLS handshake
		ResponseHeaderTimeout: 10 * time.Second,          // Set a timeout for waiting on response headers
		Proxy:                 http.ProxyFromEnvironment, // Use system-wide proxy settings
	},
}
