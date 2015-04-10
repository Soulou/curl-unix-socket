# HTTP requests over UNIX socket

I wanted to try docker API without binding it on a TCP socket.
Curl doesn't work on UNIX socket, so I've developped that:

## Build

`go get github.com/Soulou/curl-unix-socket`

### Archlinux

You can install `curl-unix-socket-git` from AUR, thank you @slopjong

## RUN

`./curl-unix-socket unix:///var/run/docker.sock:/v1.6/images/json`

### Flags

* `-X`: HTTP Verb [GET|POST|DELETE|...]
* `-d`: Request data
* `-H`: Additional Headers
  * Example: `-H 'Accept: application/json|Content-type: application/json'`
* `-b`: Add Cookie
  * Example: `-b 'Key=Value|Key2=Value2'`
* `-https`: Make an HTTPS request over the unix socket
* `-k`: Allow unsecure HTTPS requests, don't check certificate
* `-v`: Verbose
