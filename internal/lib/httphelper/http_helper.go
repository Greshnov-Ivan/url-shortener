package httphelper

import (
	"fmt"
	"net/http"
	"net/url"
)

func CreateURL(r *http.Request, path string) (string, error) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	host := r.Host
	if host == "" {
		return "", fmt.Errorf("host is empty")
	}

	baseURL := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	/*
		baseURL := url.URL{
				Host: r.Host,
				Path: path,
			}
	*/
	return baseURL.String(), nil
}
