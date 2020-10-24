package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (c *Client) ReAuthenticate() error {
	return c.Authenticate()
}

func (c *Client) Authenticate() error {
	u := c.baseURL()
	u.Path = "/api/auth"

	type requestParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	param := &requestParams{
		Username: c.c.UnifiUsername(),
		Password: c.c.UnifiPassword(),
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(param); err != nil {
		return errors.WithStack(err)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &b)
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.httpclient.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	token := res.Header.Get("Authorization")
	if token == "" {
		return errors.New("could not get Authorization header from acquireCookie response header")
	}

	header := make(http.Header)
	header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	header.Add("Content-Type", "application/json")

	c.header = header

	c.logger.Debugln("logged in to unifi")
	return nil
}
