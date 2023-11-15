package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"time"
)

var log = config.Logger

func Start(ctx context.Context, router http.Handler) {
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", config.App.Port),
		Handler: panicRecovery(router),
	}

	go func() {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err := s.Shutdown(timeoutCtx)
		if err != nil {
			log.Error().Err(err).Msgf("Error shutting down server: %s", err)
			return
		}
		log.Info().Msg("Server stopped. Will not handle incoming requests from now")
	}()

	log.Info().Msgf("chitchat starting. Listening on %s", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info().Msg("Server has been closed.")
		} else {
			log.Fatal().Err(err).Msg("Server should have stopped, but hasn't. Exiting process now.")
		}
	}
}

// PanicRecovery is a middleware that will recover any Panic that happens during the request, and will log the stack
// trace for debugging.
func panicRecovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var requestString string
				requestDump, err := httputil.DumpRequest(r, true)
				if err != nil {
					requestString = "Unable to dump request"
				} else {
					requestString = string(requestDump)
				}
				log.Error().Err(err).Any("stack_trace", debug.Stack()).Msgf("Panic while serving %s", requestString)
				app.Redirect(w, r, "/errors/internal-server-error")
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
