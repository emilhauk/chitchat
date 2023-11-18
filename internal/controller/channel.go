package controller

import (
	"errors"
	"fmt"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetChannel(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
	channelUUID := chi.URLParam(r, "channelUUID")
	channel, err := chatService.Get(channelUUID, user)

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
		channels, listErr := channelManager.GetChannelListForUser(user.UUID)
		data := map[string]any{
			"User":     user,
			"Channels": channels,
			"Channel":  channel,
		}
		if listErr != nil {
			_ = templates.ExecuteTemplate(w, "error-page", map[string]any{"Code": 500})
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
		_ = templates.ExecuteTemplate(w, "chat", data)
	}
}
