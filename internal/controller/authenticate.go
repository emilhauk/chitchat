package controller

import (
	app "github.com/emilhauk/chitchat/internal"
	internalMiddleware "github.com/emilhauk/chitchat/internal/middleware"
	"net/http"
	"strings"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.Redirect(w, r, "/")
		return
	}
	email := strings.ToLower(r.FormValue("email"))
	plainPassword := r.FormValue("password")

	user, err := userManager.FindByEmailAndPlainPassword(email, plainPassword)
	if err != nil {
		app.Redirect(w, r, "/")
		return
	}

	session, err := sessionManager.CreateSession(user.UUID)
	if err != nil {
		app.Redirect(w, r, "/")
		return
	}

	cookie := http.Cookie{
		Name:     internalMiddleware.AuthCookie,
		Value:    session.ID,
		Path:     "/",
		Domain:   r.URL.Host,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	app.Redirect(w, r, "/im")
}
