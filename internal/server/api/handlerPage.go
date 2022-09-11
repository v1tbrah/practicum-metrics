package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (a *api) handlerGetPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.handlerGetPage started")
		defer log.Debug().Msg("api.handlerGetPage ended")

		pageData := newDataForPage()
		dataMetrics, err := a.service.GetData(r.Context())
		if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}
		fillMetricsForPage(&pageData.Data, dataMetrics)
		page, err := pageData.page()
		if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		a.respond(w, http.StatusOK, []byte(page))
	}
}
