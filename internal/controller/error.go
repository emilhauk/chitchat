package controller

import (
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "internal-server-error", map[string]any{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render internal-server-error template properly")
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "bad-request", map[string]any{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render bad-request template properly")
	}
}
