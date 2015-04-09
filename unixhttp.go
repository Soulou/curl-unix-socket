package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	method, data, Cookie, Header  string
	Verbose, Https, UnsecureHttps bool
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
	flag.BoolVar(&Https, "https", false, "Make an HTTPS request")
	flag.BoolVar(&UnsecureHttps, "k", false, "Accept any certificate")
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

	hostAndPath := strings.SplitN(u.Path, ":", 2)
	if len(hostAndPath) < 2 {
		usage()
		os.Exit(1)
	}

	u.Host = hostAndPath[0]
	u.Path = hostAndPath[1]

	reader := strings.NewReader(data)
	if len(data) > 0 {
		// If there are data the request can't be GET (curl behavior)
		if method == "GET" {
			method = "POST"
		}
		// If data begins with @, it references a file
		if string(data[0]) == "@" && len(data) > 1 {
			if string(data[1:]) == "-" {
				buf, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					fmt.Println("Failed to read from stdin:", err)
					os.Exit(1)
				}
				reader = strings.NewReader(string(buf))
			} else {
				buf, err := ioutil.ReadFile(string(data[1:]))
				if err != nil {
					fmt.Println("Failed to open file:", err)
					os.Exit(1)
				}
				reader = strings.NewReader(string(buf))
			}
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

	var conn net.Conn
	if Https {
		config := &tls.Config{}
		if UnsecureHttps {
			config.InsecureSkipVerify = true
		}
		conn, err = tls.Dial("unix", u.Host, config)
	} else {
		conn, err = net.Dial("unix", u.Host)
	}

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
	defer res.Body.Close()

	if Verbose {
		fmt.Println(">")
		fmt.Printf("< %v %v\n", res.Proto, res.Status)
		for name, value := range res.Header {
			fmt.Printf("< %v: %v\n", name, strings.Join(value, " "))
		}
	}

	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil && err != io.EOF {
		fmt.Println("Invalid body in answer", err)
		os.Exit(1)
	}

	fmt.Println()
}
