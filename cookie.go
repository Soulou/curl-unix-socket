package main

import (
	"fmt"
	"net/http"
	"strings"
)

func addCookies(req *http.Request) error {
	cookies, err := parseCookies()
	if err != nil {
		return err
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	return nil
}

func parseCookies() ([]*http.Cookie, error) {
	cookies := []*http.Cookie{}
	if len(Cookie) == 0 {
		return nil, nil
	}
	cookiesStr := strings.Split(Cookie, "|")
	for _, cookieStr := range cookiesStr {
		cookieArr := strings.Split(cookieStr, "=")
		if len(cookieArr) != 2 {
			return nil, fmt.Errorf("Invalid cookie:", cookieStr)
		}
		cookies = append(cookies, &http.Cookie{
			Name:  strings.TrimSpace(cookieArr[0]),
			Value: strings.TrimSpace(cookieArr[1]),
		})
	}
	return cookies, nil
}
