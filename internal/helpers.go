package app

import (
	"context"
	"github.com/emilhauk/chitchat/internal/model"
	"net/http"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var UserContextKey = contextKey("user")

func GetUserFromContextOrPanic(ctx context.Context) model.User {
	// This will panic and that's ok. Should never be used if un certain.
	return ctx.Value(UserContextKey).(model.User)
}

func ContextWithUser(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
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
