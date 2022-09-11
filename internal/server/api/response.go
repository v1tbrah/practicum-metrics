package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (a *api) respond(w http.ResponseWriter, code int, data []byte) {
	w.WriteHeader(code)
	w.Write(data)
}

func (a *api) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	log.Error().Err(err).Msg(r.RequestURI)
	http.Error(w, err.Error(), code)
}
