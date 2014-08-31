package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func requestExecute(conn net.Conn, client *httputil.ClientConn, req *http.Request) (*http.Response, error) {
	if Verbose {
		fmt.Printf("> %s %s %s\n", req.Method, req.URL.Path, req.Proto)
		fmt.Printf("> Socket: %s\n", conn.RemoteAddr())
		for k, v := range req.Header {
			fmt.Printf("> %s: %s\n", k, strings.Join(v, " "))
		}
		fmt.Println("> Content-Length:", req.ContentLength)
	}
	return client.Do(req)
}
