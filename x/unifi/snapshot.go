package unifi

import (
	"context"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"
)

func (c *Client) GetSnapshot(ctx context.Context, doorbellID string) ([]byte, error) {
	u := c.baseURL()
	u.Path = "/api/cameras/" + doorbellID + "/snapshot"

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelFunc()

	res, err := c.request(timeoutCtx, http.MethodGet, u, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to get snapshot: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Errorln(err)
		}
	}()
	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("failed to read image at once: %w", err)
	}

	return image, nil
}
