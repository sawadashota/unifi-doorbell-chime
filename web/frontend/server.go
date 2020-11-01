package frontend

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

//go:embed static
var static embed.FS

type staticFS struct {
	s embed.FS
}

func newStaticFS() *staticFS {
	return &staticFS{
		s: static,
	}
}

func (f *staticFS) Open(name string) (fs.File, error) {
	return f.s.Open(f.resolvePath(name))
}

func (f *staticFS) ReadFile(name string) ([]byte, error) {
	return f.s.ReadFile(f.resolvePath(name))
}

func (f *staticFS) exist(name string) bool {
	file, err := f.s.Open(f.resolvePath(name))
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}

func (f *staticFS) resolvePath(name string) string {
	return filepath.Join("static", name)
}

type Server struct {
	r      Registry
	c      Configuration
	logger logrus.FieldLogger
	static *staticFS
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
		logger: r.AppLogger("frontend"),
		static: newStaticFS(),
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

	if s.static.exist(r.URL.Path) {
		http.FileServer(http.FS(s.static)).ServeHTTP(w, r)
		return
	}

	if b, err := s.static.ReadFile("index.html"); err == nil {
		w.Header().Add("Content-Type", "text/html")
		_, _ = w.Write(b)
		return
	}

	s.logger.Errorf("cannot serve %s", r.URL.Path)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) Start(ctx context.Context) error {
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.c.WebPort()),
		Handler: s,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Infof("start frontend server. 127.0.0.1:%d", s.c.WebPort())
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
		return xerrors.Errorf("exit frontend server: %w", err)
	}
}
