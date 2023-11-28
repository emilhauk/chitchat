package controller

import (
	"bytes"
	"fmt"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/emilhauk/chitchat/internal/sse"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
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

	go func(message model.Message) {
		message.Direction = model.DirectionIn
		buf := bytes.Buffer{}
		err = templates.ExecuteTemplate(&buf, "message", message)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to execute template")
			return
		}
		err := sse.PublishUsingBrokerInContext(r.Context(), sse.NewEvent("message", channelUUID, message.Sender.UUID, buf.String()))
		if err != nil {
			log.Error().Err(err).Msgf("Failed to publish event")
		}
	}(message)

	if app.IsHtmxRequest(r) {
		err = templates.ExecuteTemplate(w, "message", message)
	} else {
		app.Redirect(w, r, fmt.Sprintf("/im/channel/%s", channelUUID))
	}
}
