package controller

import (
	"errors"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetChannel(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	channelUUID := chi.URLParam(r, "channelUUID")
	channel, err := channelManager.GetChannelForUser(channelUUID, user.UUID)

	if app.IsHtmxRequest(r) {
		if err != nil {
			switch {
			case errors.Is(err, app.ErrChannelNotFound):
				_ = templates.ExecuteTemplate(w, "error-main", map[string]any{"Code": 404, "Message": "Chat not found."})
				return
			default:
				log.Error().Err(err).Msgf("Failed to load channel=%s for user=%s", channelUUID, user.UUID)
				_ = templates.ExecuteTemplate(w, "error-main", map[string]any{"Code": 500})
			}
		}
	} else {
		channels, listErr := channelManager.GetChannelListForUser(user.UUID)
		data := map[string]any{
			"User":       user,
			"Channels":   channels,
			"GetChannel": channel,
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
