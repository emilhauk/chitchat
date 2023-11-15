package middleware

import (
	"net/http"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().
			Ctx(r.Context()).
			Msgf("%s %s: %s", r.Proto, r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}
