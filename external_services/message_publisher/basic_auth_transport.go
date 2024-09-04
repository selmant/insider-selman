package message_publisher

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type BasicAuthTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
	BaseURL   string
}

func (c *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	auth := c.Username + ":" + c.Password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	if !strings.HasPrefix(req.URL.String(), "http") {
		fullURL, err := url.JoinPath(c.BaseURL, req.URL.Path)
		if err != nil {
			return nil, err
		}
		req.URL, err = url.Parse(fullURL)
		if err != nil {
			return nil, err
		}
	}

	return c.Transport.RoundTrip(req)
}

func NewBasicAuthTransport() *BasicAuthTransport {
	return &BasicAuthTransport{
		Transport: http.DefaultTransport,
		Username:  os.Getenv("WEBHOOK_USERNAME"),
		Password:  os.Getenv("WEBHOOK_PASSWORD"),
		BaseURL:   os.Getenv("WEBHOOK_BASE_URL"),
	}
}
