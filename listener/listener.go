package listener

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/browser"
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
	WebPort() uint64
}

func New(r Registry, c Configuration) Strategy {
	return &PollingStrategy{
		r:      r,
		c:      c,
		logger: r.AppLogger("listener"),
	}
}

func (s *PollingStrategy) poll(ctx context.Context) error {
	ds, err := s.r.UnifiClient().GetDoorbells(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, d := range ds {
		if d.DoesRung(s.state) {
			if err := s.openURL(&d); err != nil {
				return errors.WithStack(err)
			}

			s.logger.Infof("%s (%s) is rung!\n", d.Name, d.Mac)
		}
	}

	s.state = ds
	return nil
}

const pollingInterval = 1 * time.Second

func (s *PollingStrategy) Start() error {
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	defer func() {
		s.cancelFunc()
		s.logger.Info("Bye!")
		s.isComplete = true
	}()

	if err := s.r.UnifiClient().Authenticate(); err != nil {
		s.logger.Error(err)
		return errors.WithStack(err)
	}

	doorbells, err := s.r.UnifiClient().GetDoorbells(s.ctx)
	if err != nil {
		s.logger.Error(err)
		return errors.WithStack(err)
	}

	for _, d := range doorbells {
		s.logger.Infof("activate %s ID: %s\n", d.Name, d.ID)
	}

	ticker := time.NewTicker(pollingInterval)

	for {
		select {
		case <-s.ctx.Done():
			return nil

		case <-ticker.C:
			if err := s.poll(s.ctx); err != nil {
				s.logger.Error(err)
				return errors.WithStack(err)
			}
		}
	}
}

var shutdownPollInterval = 500 * time.Millisecond

func (s *PollingStrategy) Shutdown(ctx context.Context) error {
	ticker := time.NewTicker(shutdownPollInterval)

	defer ticker.Stop()

	s.cancelFunc()
	for {
		if s.isComplete {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *PollingStrategy) openURL(doorbell *unifi.Doorbell) error {
	return browser.OpenURL(
		fmt.Sprintf("http://127.0.0.1:%d/ringing/%s", s.c.WebPort(), doorbell.ID),
	)
}
