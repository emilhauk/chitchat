package middleware

import (
	"errors"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"net/http"
	"time"
)

const AuthCookie = "chitchat-session"

var (
	ErrCookieNotFound = errors.New("cookie not found")
)

type UserManager interface {
	FindByUUID(uuid string) (model.User, error)
}

type SessionManager interface {
	FindByID(id string) (model.Session, error)
	SetLastSeenAt(id string, lastSeenAt time.Time) error
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
		if err != nil {
			switch {
			case errors.Is(err, ErrCookieNotFound):
				app.Redirect(w, r, "/")
			case errors.Is(err, app.ErrSessionNotFound):
				app.Redirect(w, r, "/?reason=invalid-session")
			case errors.Is(err, app.ErrUserNotFound):
				app.Redirect(w, r, "/?reason=user-deleted")
			default:
				log.Error().Err(err).Msg("Unknown error establishing user login state")
				app.Redirect(w, r, "/error/internal-server-error")
			}
			log.Debug().Err(err).Msg("Authorization failed.")
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
				case errors.Is(err, ErrCookieNotFound):
					fallthrough
				case errors.Is(err, app.ErrSessionNotFound):
					fallthrough
				case errors.Is(err, app.ErrUserNotFound):
					log.Debug().Err(err).Any("cookie", r.Header.Get("cookie")).Msg("Not logged in")
					next.ServeHTTP(w, r) // It's OK for the user not to be logged in
				default:
					log.Error().Err(err).Msg("Unknown error establishing user login state")
					app.Redirect(w, r, "/error/internal-server-error")
				}
				return
			}

			log.Debug().Msgf("Hit guest pages as logged-in user. Redirect=%s", location)
			app.Redirect(w, r, location)
		})
	}
}

func (m Auth) resolveUser(r *http.Request) (model.User, error) {
	var user model.User

	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return user, ErrCookieNotFound
	}

	session, err := m.sessionManager.FindByID(cookie.Value)
	if err != nil {
		return user, err
	}
	if session.LastSeenAt == nil || session.LastSeenAt.Before(time.Now().Add(1*time.Hour)) {
		err = m.sessionManager.SetLastSeenAt(session.ID, time.Now())
		if err != nil {
			log.Error().Err(err).Msgf("Failed to update last seen for session=%s", session.ID)
		}
	}

	user, err = m.userManager.FindByUUID(session.UserUUID)
	if err != nil {
		return user, err
	}
	if user.DeactivatedAt != nil {
		return user, app.ErrUserDeactivated
	}

	return user, nil
}
