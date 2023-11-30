package middleware

import (
	"errors"
	"fmt"
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	AuthCookie        = "chitchat-session"
	RequestedURLParam = "requested-url"
)

type UserManager interface {
	FindByUUID(uuid string) (model.User, error)
}

type SessionManager interface {
	FindByIdAndMarkSeen(id string) (model.Session, error)
}

type Auth struct {
	userManager    UserManager
	sessionManager SessionManager
}

func NewAuthMiddleware(userManager UserManager, sessionManager SessionManager) Auth {
	return Auth{
		userManager:    userManager,
		sessionManager: sessionManager,
	}
}

func (m Auth) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := m.resolveUser(r)
		params := url.Values{
			RequestedURLParam: []string{fmt.Sprintf("%s%s", config.App.PublicURL, r.URL)},
		}
		if err != nil {
			redirectURL := url.URL{Path: "/"}
			switch {
			case errors.Is(err, http.ErrNoCookie):
				fallthrough
			case errors.Is(err, app.ErrSessionNotFound):
				fallthrough
			case errors.Is(err, app.ErrUserNotFound):
				log.Debug().Err(err).Msg("Authorization failed.")
			default:
				log.Error().Err(err).Msg("Unhandled error establishing user login state.")
				redirectURL.Path = "/error/internal-server-error"
			}
			redirectURL.RawQuery = params.Encode()
			app.Redirect(w, r, redirectURL.String())
			return
		}
		log.Debug().Any("user_uuid", user.UUID).Msg("User logged in.")
		next.ServeHTTP(w, r.WithContext(app.ContextWithUser(r.Context(), user)))
	})
}

func (m Auth) RedirectIfLoggedIn(location string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := m.resolveUser(r)
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					fallthrough
				case errors.Is(err, app.ErrSessionNotFound):
					DeleteSessionCookie(w, r)
					next.ServeHTTP(w, r) // It's OK for the user not to be logged in
				default:
					log.Error().Err(err).Msg("Unknown error establishing user login state.")
					app.Redirect(w, r, "/error/internal-server-error")
				}
				return
			}

			log.Debug().Msgf("Hit guest pages as logged-in user. Redirect=%s.", location)
			app.Redirect(w, r, location)
		})
	}
}

func (m Auth) resolveUser(r *http.Request) (user model.User, err error) {
	id, err := GetSessionID(r)
	if err != nil {
		return user, err
	}

	session, err := m.sessionManager.FindByIdAndMarkSeen(id)
	if err != nil {
		return user, err
	}

	user, err = m.userManager.FindByUUID(session.UserUUID)
	if err != nil {
		return user, err
	}
	if user.DeactivatedAt != nil {
		return user, app.ErrUserDeactivated
	}

	return user, err
}

func GetSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return "", err
	}

	return cookie.Value, err
}

func SetSessionCookie(w http.ResponseWriter, r *http.Request, session model.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookie,
		Value:    session.ID,
		Path:     "/",
		Domain:   r.URL.Host,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func DeleteSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookie,
		Value:    "",
		Path:     "/",
		Domain:   r.URL.Host,
		Expires:  time.Unix(0, 0),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ExtractAllowedSearchParams(u *url.URL) url.Values {
	queryParams := url.Values{}
	if requestedUrl, err := url.Parse(u.Query().Get(RequestedURLParam)); err == nil && requestedUrl != nil {
		if strings.HasPrefix(requestedUrl.String(), config.App.PublicURL) {
			queryParams.Add(RequestedURLParam, requestedUrl.String())
		} else {
			log.Warn().Str(RequestedURLParam, requestedUrl.String()).Msgf("Ignoring invalid requested url in search params. %s", config.App.PublicURL)
		}
	}
	return queryParams
}
