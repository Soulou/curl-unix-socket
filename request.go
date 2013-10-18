package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func requestExecute(client *httputil.ClientConn, req *http.Request) (*http.Response, error) {
	if Verbose {
		for k, v := range req.Header {
			fmt.Printf("> %s: %s\n", k, v)
		}
	}


	return client.Do(req)
}
