package unifi

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func (c *Client) GetSnapshot(ctx context.Context, doorbellID string) ([]byte, error) {
	u := c.baseURL()
	u.Path = "/api/cameras/" + doorbellID + "/snapshot"

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelFunc()

	res, err := c.request(timeoutCtx, http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Errorln(err)
		}
	}()
	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return image, nil
}
