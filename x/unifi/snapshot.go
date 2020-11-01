package unifi

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

func (c *Client) GetSnapshot(ctx context.Context, w io.Writer, doorbellID string) error {
	u := c.baseURL()
	u.Path = "/api/cameras/" + doorbellID + "/snapshot"

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelFunc()

	res, err := c.request(timeoutCtx, http.MethodGet, u, nil)
	if err != nil {
		return xerrors.Errorf("failed to get snapshot: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Errorln(err)
		}
	}()

	if res.StatusCode >= 300 {
		return &HttpError{
			message: res.Status,
			code:    res.StatusCode,
			url:     u,
			method:  http.MethodGet,
		}
	}

	_, _ = io.Copy(w, res.Body)

	return nil
}
