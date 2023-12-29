package controller

import (
	"fmt"
	internalMiddleware "github.com/emilhauk/chitchat/internal/middleware"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{}
	if values := internalMiddleware.ExtractAllowedSearchParams(r.URL); len(values) > 0 {
		data["QueryString"] = fmt.Sprintf("?%s", values.Encode())
	}
	_ = tmpl.ExecuteTemplate(w, "welcome", data)
}
