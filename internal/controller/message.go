package controller

import (
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	channelUUID := chi.URLParam(r, "channelUUID")
	err := r.ParseForm()
	if err != nil {
		app.Redirect(w, r, "/error/bad-request")
		return
	}
	content := r.FormValue("message")

	channel, err := channelManager.GetChannelForUser(channelUUID, user.UUID)
	if err != nil {
		// TODO Decide how to handle this scenario. User tries to send a message to a channel which does not exist or user aren't a member
		app.Redirect(w, r, "/error/bad-request")
		return
	}
	message, err := messageManager.Send(channel, model.Message{
		Sender:  user,
		Content: content,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send message to channel=%s for user=%s", channelUUID, user.UUID)

		// TODO Decide how to handle this scenario. Sending message failed
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	if app.IsHtmxRequest(r) {
		_ = templates.ExecuteTemplate(w, "message", message)
	} else {
		app.Redirect(w, r, r.URL.String())
	}
}
