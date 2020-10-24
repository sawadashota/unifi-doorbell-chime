package unifi

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Client struct {
	c      Configuration
	r      Registry
	logger logrus.FieldLogger

	httpclient *http.Client
	header     http.Header
}

type Registry interface {
	AppLogger(app string) logrus.FieldLogger
}

type Configuration interface {
	UnifiIp() string
	UnifiUsername() string
	UnifiPassword() string
}

func NewClient(r Registry, config Configuration, httpclient *http.Client) (*Client, error) {
	c := &Client{
		c:          config,
		r:          r,
		httpclient: httpclient,
		logger:     r.AppLogger("unifi-client"),
	}

	if err := c.Authenticate(); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}

func (c *Client) baseURL() *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host:   c.c.UnifiIp() + ":7443",
	}
	return u
}
