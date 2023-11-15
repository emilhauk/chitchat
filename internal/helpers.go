package app

import (
	"context"
	"github.com/emilhauk/chitchat/internal/model"
	"net/http"
)

func GetUserFromContextOrPanic(ctx context.Context) model.User {
	// This will panic and that's ok. Should never be used if un certain.
	return ctx.Value("user").(model.User)
}

func Redirect(w http.ResponseWriter, r *http.Request, location string) {
	if IsHtmxRequest(r) {
		w.Header().Add("HX-Redirect", location)
	} else {
		w.Header().Add("Location", location)
	}
	w.WriteHeader(http.StatusFound)
}

func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("hx-request") == "true"
}
