package controller

import (
	"errors"
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	internalMiddleware "github.com/emilhauk/chitchat/internal/middleware"
	"github.com/emilhauk/chitchat/internal/model"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

func CheckUsername(w http.ResponseWriter, r *http.Request) {
	log := log.With().Any("action", "CheckUsername").Logger()
	err := r.ParseForm()
	if err != nil {
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}

	email := strings.ToLower(r.FormValue("email"))
	if _, err = mail.ParseAddress(email); err != nil {
		// TODO return validation error
		app.Redirect(w, r, "/")
		return
	}

	// TODO should suppport non-hmx
	user, err := userManager.FindByEmail(email)
	if errors.Is(err, app.ErrUserNotFound) {
		verification, err := registerService.Start(email)
		if err != nil {
			app.Redirect(w, r, "/error/internal-server-error")
			return
		}
		err = templates.ExecuteTemplate(w, "register", map[string]any{
			"RegisterSession":          verification.UUID,
			"RequireEmailVerification": config.Mail.Enabled,
			"Email":                    email,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to render registration form")
		}
		return
	}
	if err != nil {
		log.Error().Err(err).Msgf("Failed to lookup user by email=%s", email)
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	err = templates.ExecuteTemplate(w, "login", map[string]any{
		"Email": user.Email,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render login form")
	}
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

	err = createAndSetSessionCookie(w, r.URL.Host, user.UUID)
	if err != nil {
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	app.Redirect(w, r, "/im")
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

	err = createAndSetSessionCookie(w, r.URL.Host, user.UUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create session")
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}

	app.Redirect(w, r, "/im?status=register-success")
}

func createAndSetSessionCookie(w http.ResponseWriter, domain, userUUID string) error {
	session, err := sessionManager.CreateSession(userUUID)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     internalMiddleware.AuthCookie,
		Value:    session.ID,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := internalMiddleware.GetSessionID(r)
	if err != nil {
		app.Redirect(w, r, "/")
		return
	}

	err = sessionManager.Delete(sessionID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Delete user session")
	}

	cookie := http.Cookie{
		Name:     internalMiddleware.AuthCookie,
		Value:    "",
		Path:     "/",
		Domain:   r.URL.Host,
		Expires:  time.Unix(0, 0),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	app.Redirect(w, r, "/")
}
