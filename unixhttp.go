package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	method, data, Cookie, Header string
	Verbose              bool
)

func usage() {
	flag.Usage()
	fmt.Println("\n→ ./curl-unix-socket [options] <URL: unix:///path/file.sock:/path>")
}

func setupFlags() {
	flag.StringVar(&method, "X", "GET", "Method of the HTTP request")
	flag.StringVar(&data, "d", "", "Body to send in the request")
	flag.StringVar(&Header, "H", "", "Additional headers: k1:v1|k2:v2|...")
	flag.StringVar(&Cookie, "b", "", "Add cookies: k1=v1|k2=v2|...") // b because thats what curl is
	flag.BoolVar(&Verbose, "v", false, "Verbose information")
	flag.Parse()
}

func checkURL() (*url.URL, error) {
	u, err := url.Parse(flag.Args()[0])
	if err != nil {
		return nil, err
	}
	if u.Scheme != "unix" {
		return nil, fmt.Errorf("Scheme must be unix ie. unix:///var/run/daemon/sock:/path")
	}
	return u, nil
}

func main() {
	setupFlags()
	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}
	u, err := checkURL()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hostAndPath := strings.Split(u.Path, ":")
	u.Host = hostAndPath[0]
	u.Path = hostAndPath[1]

	reader := strings.NewReader(data)
	if len(data) > 0 {
		// If there are data the request can't be GET (curl behavior)
		if method == "GET" {
			method = "POST"
		}
	}

	query := ""
	if len(u.RawQuery) > 0 {
		query = "?" + u.RawQuery
	}
	req, err := http.NewRequest(method, u.Path+query, reader)
	if err != nil {
		fmt.Println("Fail to create http request", err)
		os.Exit(1)
	}
	if err := addHeaders(req); err != nil {
		fmt.Println("Fail to add headers:", err)
		os.Exit(1)
	}
	if err := addCookies(req); err != nil {
		fmt.Println("Fail to add cookies:", err)
		os.Exit(1)
	}

	conn, err := net.Dial("unix", u.Host)
	if err != nil {
		fmt.Println("Fail to connect to", u.Host, ":", err)
		os.Exit(1)
	}
	client := httputil.NewClientConn(conn, nil)
	res, err := requestExecute(conn, client, req)
	if err != nil {
		fmt.Println("Fail to achieve http request over unix socket", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Invalid body in answer")
		os.Exit(1)
	}

	fmt.Println(string(body))
	defer res.Body.Close()

}
