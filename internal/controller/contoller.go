package controller

import (
	"github.com/emilhauk/chitchat/config"
	"github.com/emilhauk/chitchat/internal/manager"
	"github.com/emilhauk/chitchat/internal/service"
	"html/template"
)

var (
	log       = config.Logger
	err       error
	templates *template.Template

	userManager    manager.User
	sessionManager manager.Session
	channelManager manager.Channel
	messageManager manager.Message
	chatService    service.Chat
)

func init() {
	templates, err = template.ParseGlob("./templates/**")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse templates")
	}
}

func ProvideManagers(um manager.User, sm manager.Session, cm manager.Channel, mm manager.Message) {
	userManager = um
	sessionManager = sm
	channelManager = cm
	messageManager = mm
}

func ProvideServices(cs service.Chat) {
	chatService = cs
}
