package controller

import (
	"errors"
	"fmt"
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	internalMiddleware "github.com/emilhauk/chitchat/internal/middleware"
	"github.com/emilhauk/chitchat/internal/model"
	"net/http"
	"net/mail"
	"strings"
)

func CheckUsername(w http.ResponseWriter, r *http.Request) {
	log := log.With().Any("action", "CheckUsername").Logger()
	err := r.ParseForm()
	if err != nil {
		app.Redirect(w, r, getRequestedUrlOrDefault(r, "/error/internal-server-error"))
		return
	}
	qs := ""
	if values := internalMiddleware.ExtractAllowedSearchParams(r.URL); len(values) > 0 {
		qs = fmt.Sprintf("?%s", values.Encode())
	}

	email := strings.ToLower(r.FormValue("email"))
	if _, err = mail.ParseAddress(email); err != nil {
		// TODO return validation error
		app.Redirect(w, r, getRequestedUrlOrDefault(r, "/"))
		return
	}

	// TODO should support non-hmx
	user, err := userManager.FindByEmail(email)
	if errors.Is(err, app.ErrUserNotFound) {
		verification, err := registerService.Start(email)
		if err != nil {
			app.Redirect(w, r, getRequestedUrlOrDefault(r, "/error/internal-server-error"))
			return
		}

		err = tmpl.ExecuteTemplate(w, "register", map[string]any{
			"RegisterSession":          verification.UUID,
			"RequireEmailVerification": config.Mail.Enabled,
			"Email":                    email,
			"QueryString":              qs,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to render registration form")
		}
		return
	}
	if err != nil {
		log.Error().Err(err).Msgf("Failed to lookup user by email=%s", email)
		app.Redirect(w, r, getRequestedUrlOrDefault(r, "/error/internal-server-error"))
		return
	}
	err = tmpl.ExecuteTemplate(w, "login", map[string]any{
		"Email":       user.Email,
		"QueryString": qs,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render login form")
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.Redirect(w, r, "/")
		return
	}
	request := model.RegisterRequest{
		VerificationUUID: r.FormValue("register-session"),
		Code:             r.FormValue("code"),
		Name:             r.FormValue("name"),
		PlainPassword:    r.FormValue("password"),
	}

	user, err := registerService.Fulfill(request)

	// Don't leak password
	request.PlainPassword = ""

	if err != nil {
		// Several validation errors may occur here. Non is handled yet, just internal server error for now
		switch {
		case errors.Is(err, app.ErrFieldVerificationNotFound):
			fallthrough
		case errors.Is(err, app.ErrFieldVerificationCodeInvalid):
			fallthrough
		default:
			log.Error().Err(err).Any("user", request).Msg("Failed to register user")
			app.Redirect(w, r, "/error/internal-server-error")
			return
		}
	}

	session, err := sessionManager.CreateSession(user.UUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create session")
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	internalMiddleware.SetSessionCookie(w, r, session)
	app.Redirect(w, r, getRequestedUrlOrDefault(r, "/im?status=register-success"))
}

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
		log.Error().Err(err).Msg("Failed to create session")
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	internalMiddleware.SetSessionCookie(w, r, session)
	app.Redirect(w, r, getRequestedUrlOrDefault(r, "/im"))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := internalMiddleware.GetSessionID(r)
	if err == nil {
		err = sessionManager.Delete(sessionID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to Delete user session")
		}
	}
	internalMiddleware.DeleteSessionCookie(w, r)

	app.Redirect(w, r, "/")
}

func getRequestedUrlOrDefault(r *http.Request, defaultUrl string) string {
	values := internalMiddleware.ExtractAllowedSearchParams(r.URL)
	if requestedUrl := values.Get(internalMiddleware.RequestedURLParam); requestedUrl != "" {
		return requestedUrl
	}
	return defaultUrl
}
