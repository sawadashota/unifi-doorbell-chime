package unifi

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/xerrors"
)

func (c *Client) SetMessage(ctx context.Context, doorbellID string, message string, duration time.Duration) error {
	u := c.baseURL()
	u.Path = "/api/cameras/" + doorbellID

	type requestParams struct {
		LcdMessage struct {
			Duration uint64 `json:"duration"`
			Text     string `json:"text"`
			Type     string `json:"type"`
		} `json:"lcdMessage"`
	}

	param := &requestParams{
		LcdMessage: struct {
			Duration uint64 `json:"duration"`
			Text     string `json:"text"`
			Type     string `json:"type"`
		}{
			Duration: uint64(duration.Seconds()),
			Text:     message,
			Type:     "CUSTOM_MESSAGE",
		},
	}

	if err := c.jsonRequest(ctx, http.MethodPatch, u, param, nil); err != nil {
		return xerrors.Errorf("failed to set message: %w", err)
	}
	c.logger.Debugf(`set "%s" as message successfully`, message)

	return nil
}
