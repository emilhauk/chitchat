package controller

import (
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{}
	if r.URL.Query().Has("requested-url") {
		data["QueryString"] = "?" + r.URL.RawQuery
	}
	_ = templates.ExecuteTemplate(w, "welcome", data)
}
