package controller

import (
	"fmt"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Join(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContextOrPanic(r.Context())
	inviationCode := chi.URLParam(r, "invitationCode")

	err := chatService.AcceptInvitation(inviationCode, user.UUID)
	if err != nil {
		app.Redirect(w, r, "/error/internal-server-error")
		return
	}
	app.Redirect(w, r, fmt.Sprintf("/im/channel/%s", inviationCode))
}
