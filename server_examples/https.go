package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: %s <socket path> <cert> <key>", os.Args[0])
		os.Exit(-1)
	}
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: os.Args[1]})
	if err != nil {
		panic(err)
	}
	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(os.Args[2], os.Args[3])
	if err != nil {
		panic(err)
	}
	tlsListener := tls.NewListener(listener, config)

	log.Println(http.Serve(tlsListener, nil))
}
