package server

import (
	"github.com/emilhauk/chitchat/internal/controller"
	internalMiddleware "github.com/emilhauk/chitchat/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewRouter(authMiddleware internalMiddleware.Auth) http.Handler {
	r := chi.NewRouter()
	r.Use(internalMiddleware.RequestLogger)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RedirectIfLoggedIn("/im"))

		r.Get("/", controller.Welcome)
		r.Post("/login", controller.Login)
	})

	r.Route("/im", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuthenticatedUser)

		r.Get("/", controller.Main)

		r.Route("/channel", func(r chi.Router) {
			r.Route("/{channelUUID}", func(r chi.Router) {
				r.Get("/", controller.GetChannel)
				r.Post("/message", controller.SendMessage)
			})
		})
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	fileServer(r, "/public", filesDir)

	log.Info().Msg("chitchat starting. Listening on port 3333")
	err := http.ListenAndServe(":3333", r)
	if err != nil {
		log.Error().Err(err).Send()
	}

	return r
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
		pathPrefix := strings.TrimSuffix(chi.RouteContext(r.Context()).RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
