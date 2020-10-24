package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type HttpError struct {
	message string
	code    int
	url     *url.URL
	method  string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("message: %s. code: %d. url: %s", e.message, e.code, e.url.String())
}

func (c *Client) request(ctx context.Context, method string, u *url.URL, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header = c.header

	res, err := func() (*http.Response, error) {
		req.WithContext(ctx)
		respCh := make(chan *http.Response)
		errCh := make(chan error)

		go func() {
			resp, err := c.httpclient.Do(req)
			if err != nil {
				errCh <- err
				return
			}

			respCh <- resp
		}()

		select {
		case resp := <-respCh:
			return resp, nil

		case err := <-errCh:
			return nil, err

		case <-ctx.Done():
			return nil, errors.New("HTTP jsonRequest cancelled")
		}
	}()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (c *Client) jsonRequest(ctx context.Context, method string, u *url.URL, param, response interface{}) error {
	buf := new(bytes.Buffer)
	if param != nil {
		if err := json.NewEncoder(buf).Encode(param); err != nil {
			return errors.WithStack(err)
		}
	}

	res, err := c.request(ctx, method, u, buf)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Errorln(err)
		}
	}()

	if res.StatusCode >= 300 {
		err := &HttpError{
			message: "failed to request",
			code:    res.StatusCode,
			url:     u,
			method:  method,
		}
		c.logger.Warnf(err.Error())
		body, _ := ioutil.ReadAll(res.Body)
		c.logger.Debugln(body)
		return err
	}

	if response == nil {
		return nil
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
