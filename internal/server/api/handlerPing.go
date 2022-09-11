package api

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

func (a *api) handlerPing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.handlerPing started")
		defer log.Debug().Msg("api.handlerPing ended")

		err := a.service.PingDBStorage(r.Context())
		if err != nil {
			if err == service.ErrNotDBStorage {
				a.error(w, r, http.StatusBadRequest, err)
			} else {
				a.error(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		a.respond(w, http.StatusOK, []byte{})
	}
}
