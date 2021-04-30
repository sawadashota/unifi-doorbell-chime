package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type Server struct {
	r      Registry
	c      Configuration
	logger logrus.FieldLogger
}

type Registry interface {
	AppLogger(app string) logrus.FieldLogger
	UnifiClient() *unifi.Client
}

type Configuration interface {
	APIPort() int
	MessageList() []string
}

func New(r Registry, c Configuration) *Server {
	return &Server{
		r:      r,
		c:      c,
		logger: r.AppLogger("api"),
	}
}

func (s *Server) Start(ctx context.Context) error {
	m := mux.NewRouter()
	m.Use(s.allowCORS)
	m.Use(s.requestLogging)
	m.HandleFunc("/snapshot/{doorbellID}", s.getSnapshot).Methods(http.MethodGet)
	m.HandleFunc("/message/set", s.setMessage).Methods(http.MethodPost)
	m.HandleFunc("/message/templates", s.messageTemplateList).Methods(http.MethodGet)
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.c.APIPort()),
		Handler: m,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Infof("start API server. 127.0.0.1:%d", s.c.APIPort())
		if err := svr.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Bye!")
		return svr.Shutdown(ctx)
	case err := <-errCh:
		s.logger.Debugf("%+v", err)
		return xerrors.Errorf("exit api server: %w", err)
	}
}
