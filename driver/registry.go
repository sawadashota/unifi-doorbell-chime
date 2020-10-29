package driver

import (
	"crypto/tls"
	"net/http"
	"time"

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
	fs *frontend.Server
	as *api.Server
}

var _ Registry = new(DefaultRegistry)

func NewDefaultRegistry(config configuration.Provider) Registry {
	return &DefaultRegistry{
		c: config,
	}
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
	if d.uc == nil {
		httpclient := http.DefaultClient

		if d.c.UnifiSkipTLSVerify() {
			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			httpclient = &http.Client{
				Transport: transport,
			}
		}

		d.uc = unifi.NewClient(d, d.c, httpclient)
	}
	return d.uc
}

func (d *DefaultRegistry) Listener() listener.Strategy {
	if d.ls == nil {
		d.ls = listener.New(d, d.c)
	}
	return d.ls
}

func (d *DefaultRegistry) WebFrontendServer() *frontend.Server {
	if d.fs == nil {
		d.fs = frontend.New(d, d.c)
	}
	return d.fs
}

func (d *DefaultRegistry) WebApiServer() *api.Server {
	if d.as == nil {
		d.as = api.New(d, d.c)
	}
	return d.as
}
