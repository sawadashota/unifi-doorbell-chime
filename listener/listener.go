package listener

import (
	"context"
	"fmt"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/pkg/errors"
	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"
	"github.com/sirupsen/logrus"
)

type Strategy interface {
	Start() error
	Shutdown(ctx context.Context) error
}

var _ Strategy = new(PollingStrategy)

type PollingStrategy struct {
	state      unifi.Doorbells
	r          Registry
	c          Configuration
	ctx        context.Context
	cancelFunc context.CancelFunc
	isComplete bool
	logger     logrus.FieldLogger
}

type Registry interface {
	AppLogger(app string) logrus.FieldLogger
	UnifiClient() *unifi.Client
}

type Configuration interface {
	NotificationIcon() string
	WebPort() uint64
}

func New(r Registry, c Configuration) Strategy {
	return &PollingStrategy{
		r:      r,
		c:      c,
		logger: r.AppLogger("listener"),
	}
}

func (h *PollingStrategy) poll(ctx context.Context) error {
	ds, err := h.r.UnifiClient().GetDoorbells(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, d := range ds {
		if d.DoesRung(h.state) {
			_ = h.r.UnifiClient().SetMessage(ctx, d.ID, "呼び出し中・・・", 30*time.Second)
			if err := h.notify(&d); err != nil {
				return errors.WithStack(err)
			}

			h.logger.Infof("%s (%s) is rung!\n", d.Name, d.Mac)
		}
	}

	h.state = ds
	return nil
}

const pollingInterval = 1 * time.Second

func (h *PollingStrategy) Start() error {
	h.ctx, h.cancelFunc = context.WithCancel(context.Background())
	defer func() {
		h.cancelFunc()
		h.logger.Info("Bye!")
		h.isComplete = true
	}()

	if err := h.r.UnifiClient().Authenticate(); err != nil {
		h.logger.Error(err)
		return errors.WithStack(err)
	}

	doorbells, err := h.r.UnifiClient().GetDoorbells(h.ctx)
	if err != nil {
		h.logger.Error(err)
		return errors.WithStack(err)
	}

	for _, d := range doorbells {
		h.logger.Infof("activate %s ID: %s\n", d.Name, d.ID)
	}

	ticker := time.NewTicker(pollingInterval)

	for {
		select {
		case <-h.ctx.Done():
			return nil

		case <-ticker.C:
			if err := h.poll(h.ctx); err != nil {
				h.logger.Error(err)
				return errors.WithStack(err)
			}
		}
	}
}

var shutdownPollInterval = 500 * time.Millisecond

func (h *PollingStrategy) Shutdown(ctx context.Context) error {
	ticker := time.NewTicker(shutdownPollInterval)

	defer ticker.Stop()

	h.cancelFunc()
	for {
		if h.isComplete {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (h *PollingStrategy) notify(doorbell *unifi.Doorbell) error {
	note := gosxnotifier.NewNotification("Someone at the door")
	note.Title = doorbell.Name
	note.Sound = gosxnotifier.Glass
	note.Link = fmt.Sprintf("http://127.0.0.1:%d/ringing/%s", h.c.WebPort(), doorbell.ID)
	if h.c.NotificationIcon() != "" {
		note.AppIcon = h.c.NotificationIcon()
	}
	return note.Push()
}
