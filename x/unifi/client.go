package unifi

import (
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type Client struct {
	c      Configuration
	r      Registry
	logger logrus.FieldLogger

	httpclient          *http.Client
	authenticatedHeader http.Header
}

type Registry interface {
	AppLogger(app string) logrus.FieldLogger
}

type Configuration interface {
	UnifiIp() string
	UnifiUsername() string
	UnifiPassword() string
}

func NewClient(r Registry, config Configuration, httpclient *http.Client) *Client {
	return &Client{
		c:          config,
		r:          r,
		httpclient: httpclient,
		logger:     r.AppLogger("unifi-client"),
	}
}

func (c *Client) baseURL() *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host:   c.c.UnifiIp() + ":7443",
	}
	return u
}
