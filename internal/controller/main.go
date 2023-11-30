package controller

import (
	app "github.com/emilhauk/chitchat/internal"
	"net/http"
)

func Main(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
	channels, err := chatService.GetChannelList(user)
	data := map[string]any{
		"User":     user,
		"Channels": channels,
	}
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get channel list for user=%s", user.UUID)
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	_ = templates.ExecuteTemplate(w, "chat", data)
}
