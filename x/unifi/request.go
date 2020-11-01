package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/xerrors"
)

type HttpError struct {
	message string
	code    int
	url     *url.URL
	method  string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("message: %s. code: %d. url: %s", e.message, e.code, e.url.String())
}

func (e HttpError) Is(target error) bool {
	_, ok := target.(*HttpError)
	return ok
}

func (e HttpError) As(target interface{}) bool {
	if cast, ok := target.(**HttpError); ok {
		(*cast).message = e.message
		return true
	}
	return false
}

func (e HttpError) Code() int {
	return e.code
}

const defaultRequestTimeout = 3 * time.Second

func (c *Client) request(ctx context.Context, method string, u *url.URL, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, xerrors.Errorf("failed to create new request instance: %w", err)
	}
	req.Header = c.authenticatedHeader

	res, err := func() (*http.Response, error) {
		req = req.WithContext(ctx)
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
			return nil, xerrors.New("HTTP request cancelled")
		}
	}()
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		c.authenticatedHeader = nil

		c.logger.Info("try re authentication")
		if err := c.Authenticate(); err == nil {
			c.logger.Info("re authenticated successfully")
			return c.request(ctx, method, u, body)
		}
		c.logger.Error("failed to re authenticated")
	}

	return res, nil
}

func (c *Client) jsonRequest(ctx context.Context, method string, u *url.URL, param, response interface{}) error {
	buf := new(bytes.Buffer)
	if param != nil {
		if err := json.NewEncoder(buf).Encode(param); err != nil {
			return xerrors.Errorf("failed encode request json to %T: %w", param, err)
		}
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelFunc()

	res, err := c.request(timeoutCtx, method, u, buf)
	if err != nil {
		return xerrors.Errorf(": %w", err)
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
		c.logger.Warn(err)
		return xerrors.Errorf("failed to json request: %w", err)
	}

	if response == nil {
		return nil
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return xerrors.Errorf("failed to decode response json to %T: %w", response, err)
	}

	return nil
}
