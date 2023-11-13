package main

import (
	"context"
	"github.com/emilhauk/chitchat/internal/manager"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	templates      *template.Template
	channelManager = manager.Channel{}
)

func main() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	var err error
	templates, err = template.ParseGlob("./templates/**")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse templates")
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	fileServer(r, "/public", filesDir)

	r.Get("/", requireLoggedInUser(getIndex))
	r.Get("/channel/{uuid}", requireLoggedInUser(getChannel))
	r.Post("/channel/{uuid}", requireLoggedInUser(sendMessage))

	err = http.ListenAndServe(":3333", r)
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func requireLoggedInUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r.WithContext(context.WithValue(r.Context(), "user", model.User{UUID: "0-1-3", Name: "Emil Haukeland"})))
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	channels, err := channelManager.GetChannelsForUser(user.UUID)
	err = templates.ExecuteTemplate(w, "chat", map[string]any{
		"User":     user,
		"Channels": channels,
	})
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func getChannel(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	channelUUID := chi.URLParam(r, "uuid")
	channel, err := channelManager.GetChannelForUser(channelUUID, user.UUID)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	if isHtmxRequest(r) {
		err := templates.ExecuteTemplate(w, "channel", channel)
		if err != nil {
			log.Error().Err(err).Send()
		}
	} else {
		channels, err := channelManager.GetChannelsForUser(user.UUID)
		err = templates.ExecuteTemplate(w, "chat", map[string]any{
			"User":     r.Context().Value("user"),
			"Channels": channels,
			"Channel":  channel,
		})
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	channelUUID := chi.URLParam(r, "uuid")
	r.ParseForm()
	content := r.FormValue("message")
	message, err := channelManager.SendMessage(channelUUID, model.Message{
		Sender:  user,
		Content: content,
	})
	if err != nil {
		log.Error().Err(err).Send()
	}
	if isHtmxRequest(r) {
		err = templates.ExecuteTemplate(w, "message", message)
		if err != nil {
			log.Error().Err(err).Send()
		}
	} else {
		w.Header().Set("Location", r.URL.String())
		w.WriteHeader(http.StatusFound)
	}
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func isHtmxRequest(r *http.Request) bool {
	return r.Header.Get("hx-request") == "true"
}
