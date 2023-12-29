package controller

import (
	"github.com/emilhauk/chitchat/config"
	"github.com/emilhauk/chitchat/internal/manager"
	"github.com/emilhauk/chitchat/internal/service"
	"github.com/emilhauk/chitchat/templates"
)

var (
	log = config.Logger

	tmpl = templates.Templates

	userManager     manager.User
	sessionManager  manager.Session
	channelManager  manager.Channel
	messageManager  manager.Message
	chatService     service.Chat
	registerService service.Register
)

func ProvideManagers(um manager.User, sm manager.Session, cm manager.Channel, mm manager.Message) {
	userManager = um
	sessionManager = sm
	channelManager = cm
	messageManager = mm
}

func ProvideServices(cs service.Chat, rs service.Register) {
	chatService = cs
	registerService = rs
}
