package controller

import (
	"errors"
	"fmt"
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetChannel(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
	channelUUID := chi.URLParam(r, "channelUUID")
	channel, err := chatService.GetChannel(channelUUID, user)

	if channel.IsCurrentUserAdmin {
		channel.InvitationURL = fmt.Sprintf("%s/join/%s", config.App.PublicURL, channel.UUID)
	}

	if app.IsHtmxRequest(r) {
		if err != nil {
			switch {
			case errors.Is(err, app.ErrChannelNotFound):
				_ = templates.ExecuteTemplate(w, "error-main", map[string]any{"Code": 404, "Message": "Chat not found."})
			default:
				log.Error().Err(err).Msgf("Failed to load channel=%s for user=%s", channelUUID, user.UUID)
				_ = templates.ExecuteTemplate(w, "error-main", map[string]any{"Code": 500})
			}
			return
		}
		_ = templates.ExecuteTemplate(w, "channel", channel)
	} else {
		channels, listErr := chatService.GetChannelList(user)
		data := map[string]any{
			"User":     user,
			"Channels": channels,
			"Channel":  channel,
		}
		if listErr != nil {
			log.Error().Err(listErr).Msg("Failed to load channel list")
			app.Redirect(w, r, "/error/internal-server-error")
			return
		}
		if err != nil {
			switch {
			case errors.Is(err, app.ErrChannelNotFound):
				data["ErrorMain"] = map[string]any{"Code": 404, "Message": "Chat not found."}
			default:
				log.Error().Err(err).Msgf("Failed to load channel=%s for user=%s", channelUUID, user.UUID)
				data["ErrorMain"] = map[string]any{"Code": 500}
			}
		}
		err = templates.ExecuteTemplate(w, "chat", data)
		if err != nil {
			log.Warn().Err(err).Send()
		}
	}
}

func NewChannelForm(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
	if app.IsHtmxRequest(r) {
		_ = templates.ExecuteTemplate(w, "new-channel-form", map[string]any{})
	} else {
		channels, listErr := channelManager.GetChannelListForUser(user.UUID)
		data := map[string]any{
			"User":               user,
			"Channels":           channels,
			"ShowNewChannelForm": true,
		}
		if listErr != nil {
			_ = templates.ExecuteTemplate(w, "error-page", map[string]any{"Code": 500})
			return
		}
		_ = templates.ExecuteTemplate(w, "chat", data)
	}
}

func CreateNewChannel(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
	err := r.ParseForm()
	if err != nil {
		app.Redirect(w, r, "/error/bad-request")
		return
	}
	name := r.FormValue("name")
	channel, err := channelManager.Create(name, user)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create channel for user=%s", user.UUID)
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	app.Redirect(w, r, fmt.Sprintf("/im/channel/%s", channel.UUID))
}
