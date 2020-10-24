package driver

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sawadashota/unifi-doorbell-chime/driver/configuration"
	"github.com/sawadashota/unifi-doorbell-chime/listener"
	"github.com/sawadashota/unifi-doorbell-chime/web/api"
	"github.com/sawadashota/unifi-doorbell-chime/web/frontend"
	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"
	"github.com/sirupsen/logrus"
)

type Registry interface {
	Logger() logrus.FieldLogger
	AppLogger(app string) logrus.FieldLogger
	UnifiClient() *unifi.Client
	Listener() listener.Strategy
	WebFrontendServer() *frontend.Server
	WebApiServer() *api.Server
}

type DefaultRegistry struct {
	l  logrus.FieldLogger
	uc *unifi.Client
	ls listener.Strategy
	c  configuration.Provider
	ws *frontend.Server
	wa *api.Server
}

var _ Registry = new(DefaultRegistry)

func NewDefaultRegistry(config configuration.Provider) (Registry, error) {
	var err error
	httpclient := http.DefaultClient

	if config.UnifiSkipTLSVerify() {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpclient = &http.Client{
			Transport: transport,
		}
	}

	r := &DefaultRegistry{
		c: config,
	}

	r.uc, err = unifi.NewClient(r, config, httpclient)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.ls, err = listener.New(r, config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.ws = frontend.New(r, config)
	r.wa = api.New(r, config)

	return r, nil
}

func (d *DefaultRegistry) Logger() logrus.FieldLogger {
	if d.l == nil {
		l := logrus.New()
		if level, err := logrus.ParseLevel(d.c.LogLevel()); err == nil {
			l.SetLevel(level)
		} else {
			l.SetLevel(logrus.InfoLevel)
		}
		l.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
		d.l = l
	}

	return d.l
}

func (d *DefaultRegistry) AppLogger(app string) logrus.FieldLogger {
	return d.Logger().(*logrus.Logger).WithField("app", app)
}

func (d *DefaultRegistry) UnifiClient() *unifi.Client {
	return d.uc
}

func (d *DefaultRegistry) Listener() listener.Strategy {
	return d.ls
}

func (d *DefaultRegistry) WebFrontendServer() *frontend.Server {
	return d.ws
}

func (d *DefaultRegistry) WebApiServer() *api.Server {
	return d.wa
}
