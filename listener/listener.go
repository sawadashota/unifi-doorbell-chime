package listener

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type Listener struct {
	state  unifi.Doorbells
	r      Registry
	c      Configuration
	logger logrus.FieldLogger
}

type Registry interface {
	AppLogger(app string) logrus.FieldLogger
	UnifiClient() *unifi.Client
}

type Configuration interface {
	WebPort() int
}

func New(r Registry, c Configuration) *Listener {
	return &Listener{
		r:      r,
		c:      c,
		logger: r.AppLogger("listener"),
	}
}

func (l *Listener) poll(ctx context.Context) error {
	ds, err := l.r.UnifiClient().GetDoorbells(ctx)
	if err != nil {
		return xerrors.Errorf("failed to poll: %w", err)
	}

	for _, d := range ds {
		if d.DoesRung(l.state) {
			if err := l.onRung(&d); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	l.state = ds
	return nil
}

const pollingInterval = 1 * time.Second

func (l *Listener) Start(ctx context.Context) error {
	defer l.logger.Info("Bye!")

	if err := l.ping(ctx); err != nil {
		l.logger.Error(err)
		return errors.WithStack(err)
	}

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			if err := l.poll(ctx); err != nil {
				l.logger.Debugf("%+v", err)
				return xerrors.Errorf(": %w", err)
			}
		}
	}
}
func (l *Listener) ping(ctx context.Context) error {
	bc := backoff.NewExponentialBackOff()
	bc.MaxElapsedTime = time.Minute * 5
	bc.Reset()

	return backoff.Retry(func() error {
		if err := l.r.UnifiClient().Authenticate(); err != nil {
			l.logger.Error(err)
			return xerrors.Errorf("failed to authenticate: %w", err)
		}

		doorbells, err := l.r.UnifiClient().GetDoorbells(ctx)
		if err != nil {
			l.logger.Error(err)
			return xerrors.Errorf("failed to start listener: %w", err)
		}

		for _, d := range doorbells {
			l.logger.Infof("activate %s ID: %s\n", d.Name, d.ID)
		}
		return nil
	}, bc)
}

func (l *Listener) onRung(doorbell *unifi.Doorbell) error {
	err := browser.OpenURL(
		fmt.Sprintf("http://127.0.0.1:%d/ringing/%s", l.c.WebPort(), doorbell.ID),
	)
	if err != nil {
		return xerrors.Errorf("failed to open browser: %w", err)
	}

	l.logger.Infof("%s (%s) is rung!\n", doorbell.Name, doorbell.Mac)
	return nil
}
