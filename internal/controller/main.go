package controller

import (
	app "github.com/emilhauk/chitchat/internal"
	"net/http"
)

func Main(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
	channels, err := channelManager.GetChannelListForUser(user.UUID)
	data := map[string]any{
		"User":     user,
		"Channels": channels,
	}
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get channel list for user=%s", user.UUID)
		data["ErrorMain"] = map[string]any{"Code": 500, "Message": "Failed to load channel list please try again later."}
	}
	_ = templates.ExecuteTemplate(w, "chat", data)
}
