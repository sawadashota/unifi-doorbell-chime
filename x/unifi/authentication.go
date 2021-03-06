package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

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
		return xerrors.Errorf("failed to encode request param to %T: %w", param, err)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &b)
	if err != nil {
		return xerrors.Errorf("failed to create request instance: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.httpclient.Do(req)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(err)
		}
	}()

	token := res.Header.Get("Authorization")
	if token == "" {
		return xerrors.New("could not get Authorization Header from acquireCookie response authenticatedHeader")
	}

	header := make(http.Header)
	header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	header.Add("Content-Type", "application/json")

	c.authenticatedHeader = header

	c.logger.Debugln("logged in to unifi")
	return nil
}
