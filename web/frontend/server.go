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
	"strings"

	"github.com/sawadashota/unifi-doorbell-chime/x/unifi"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

var (
	//go:embed static/*
	static embed.FS

	public fs.FS
)

func init() {
	var err error
	public, err = fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
}

type staticFS struct {
	logger logrus.FieldLogger
}

func newStaticFS(r Registry) *staticFS {
	return &staticFS{
		logger: r.AppLogger("EmbedFileSystem"),
	}
}

func (f *staticFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(public, name)
}

func (f *staticFS) exist(name string) bool {
	v := strings.SplitN(name, "/", 2)
	if len(v) != 2 {
		f.logger.Errorf("invalid path: %s", name)
		return false
	}
	file, err := public.Open(v[1])
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
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
	WebPort() int
	ApiPort() int
	MessageList() []string
}

func New(r Registry, c Configuration) *Server {
	return &Server{
		r:      r,
		c:      c,
		logger: r.AppLogger("frontend"),
		static: newStaticFS(r),
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
		http.FileServer(http.FS(public)).ServeHTTP(w, r)
		return
	}

	if filepath.Ext(r.URL.Path) == "html" || filepath.Ext(r.URL.Path) == "" {
		if b, err := fs.ReadFile(public, "index.html"); err == nil {
			w.Header().Add("Content-Type", "text/html")
			_, _ = w.Write(b)
			return
		}
	}

	s.logger.Warnf("%s is not found", r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
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
