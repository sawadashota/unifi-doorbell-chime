package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"
	"github.com/sirupsen/logrus"
)

type Server struct {
	r      Registry
	c      Configuration
	logger logrus.FieldLogger
	svr    *http.Server
}

type Registry interface {
	AppLogger(app string) logrus.FieldLogger
	UnifiClient() *unifi.Client
}

type Configuration interface {
	ApiPort() uint64
	MessageList() []string
}

func New(r Registry, c Configuration) *Server {
	return &Server{
		r:      r,
		c:      c,
		logger: r.AppLogger("api"),
		svr: &http.Server{
			Addr: fmt.Sprintf(":%d", c.ApiPort()),
		},
	}
}

func (s *Server) Start() error {
	m := mux.NewRouter()
	m.Use(s.allowCORS)
	m.Use(s.requestLogging)
	m.HandleFunc("/snapshot/{doorbellID}", s.getSnapshot).Methods(http.MethodGet)
	m.HandleFunc("/message/set", s.setMessage).Methods(http.MethodPost)
	m.HandleFunc("/message/templates", s.messageTemplateList).Methods(http.MethodGet)
	s.svr.Handler = m

	s.logger.Infof("start API server. 127.0.0.1:%d", s.c.ApiPort())
	return s.svr.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Bye!")
	return s.svr.Shutdown(ctx)
}
