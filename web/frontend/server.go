package frontend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"

	"github.com/gobuffalo/packr/v2"
	"github.com/sirupsen/logrus"
)

type Server struct {
	r      Registry
	c      Configuration
	static *packr.Box
	svr    *http.Server
	logger logrus.FieldLogger
}

type Registry interface {
	AppLogger(app string) logrus.FieldLogger
	UnifiClient() *unifi.Client
}

type Configuration interface {
	WebPort() uint64
	ApiPort() uint64
	MessageList() []string
}

func New(r Registry, c Configuration) *Server {
	return &Server{
		r:      r,
		c:      c,
		static: packr.New("static files", "./build"),
		logger: r.AppLogger("frontend"),
	}
}

func (s *Server) handleWellKnownConfiguration(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		ApiEndpoint string `json:"api_endpoint"`
	}{
		ApiEndpoint: fmt.Sprintf("http://127.0.0.1:%d", s.c.ApiPort()),
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&res); err != nil {
		s.logger.Error(err)
	}
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(buf.Bytes())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/.well-known/configuration" {
		s.handleWellKnownConfiguration(w, r)
		return
	}

	if _, err := s.static.Find(r.URL.Path); err == nil {
		http.FileServer(s.static).ServeHTTP(w, r)
		return
	}

	if b, err := s.static.Find("index.html"); err == nil {
		w.Header().Add("Content-Type", "text/html")
		_, _ = w.Write(b)
		return
	}
	s.logger.Errorf("cannot serve %s", r.URL.Path)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) Start() error {
	s.logger.Infof("start frontend server. 127.0.0.1:%d", s.c.WebPort())
	s.svr = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.c.WebPort()),
		Handler: s,
	}
	if err := s.svr.ListenAndServe(); err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Bye!")
	return s.svr.Shutdown(ctx)
}
