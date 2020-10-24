package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) messageTemplateList(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Templates []string `json:"templates"`
	}{
		Templates: s.c.MessageList(),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&res); err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(buf.Bytes())
}

func (s *Server) getSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doorbellID := vars["doorbellID"]
	image, err := s.r.UnifiClient().GetSnapshot(r.Context(), doorbellID)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(image)
}

func (s *Server) setMessage(w http.ResponseWriter, r *http.Request) {
	param := struct {
		DoorbellID  string `json:"doorbell_id"`
		Message     string `json:"message"`
		DurationSec uint64 `json:"duration_sec"`
	}{
		DurationSec: 30,
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			s.logger.Error(err)
		}
	}()
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.r.UnifiClient().SetMessage(
		r.Context(),
		param.DoorbellID,
		param.Message,
		time.Duration(param.DurationSec)*time.Second,
	); err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
