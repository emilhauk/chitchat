package controller

import "net/http"

func Welcome(w http.ResponseWriter, r *http.Request) {
	_ = templates.ExecuteTemplate(w, "welcome", map[string]any{})
}
