package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
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
	RedirectWithReturnURL(w, r, location, "")
}

func RedirectWithReturnURL(w http.ResponseWriter, r *http.Request, location, returnUrl string) {
	if returnUrl != "" {
		if u, err := url.Parse(location); err != nil {
			log.Error().Err(err).Msgf("Failed to parse location=%s", location)
		} else {
			u.Query().Add("return-url", url.QueryEscape(returnUrl))
			location = u.String()
		}
	}
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

func BuildGravatar(email string) string {
	hash := hex.EncodeToString(sha256.New().Sum([]byte(strings.ToLower(strings.TrimSpace(email)))))
	return fmt.Sprintf("https://gravatar.com/avatar/%s?r=r&d=retro", hash)
}

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
