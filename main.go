package main

import (
	"context"
	"github.com/emilhauk/chitchat/config"
	"github.com/emilhauk/chitchat/internal/controller"
	"github.com/emilhauk/chitchat/internal/database"
	"github.com/emilhauk/chitchat/internal/manager"
	internalMiddleware "github.com/emilhauk/chitchat/internal/middleware"
	"github.com/emilhauk/chitchat/internal/server"
	"github.com/emilhauk/chitchat/internal/service"
	"github.com/emilhauk/chitchat/internal/sse"
)

var (
	log                 = config.Logger
	dbStore             database.DBStore
	userManager         manager.User
	sessionManager      manager.Session
	channelManager      manager.Channel
	messageManager      manager.Message
	verificationManager manager.Verification
	credentialManager   manager.Credential
	chatService         service.Chat
	registerService     service.Register
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var db, err = database.NewConnectionPool(config.Database)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer db.Close()

	dbStore = database.NewDBStore(db)
	userManager = manager.NewUserManager(dbStore.Users, dbStore.Credentials)
	sessionManager = manager.NewSessionManager(dbStore.Sessions)
	channelManager = manager.NewChannelManager(dbStore.Channels)
	messageManager = manager.NewMessageManager(dbStore.Messages)
	verificationManager = manager.NewVerificationManager(dbStore.Verifications)
	credentialManager = manager.NewCredentialManager(dbStore.Credentials)

	chatService = service.NewChatService(userManager, channelManager, messageManager)
	registerService = service.NewRegisterService(userManager, verificationManager, credentialManager)

	// TODO This stinks. Should provide better wrapper for controllers
	controller.ProvideManagers(userManager, sessionManager, channelManager, messageManager)
	controller.ProvideServices(chatService, registerService)

	authMiddleware := internalMiddleware.NewAuthMiddleware(userManager, sessionManager)
	sseBroker := sse.NewBroker(config.Logger)
	router := server.NewRouter(authMiddleware, sseBroker)

	server.Start(ctx, router)
}
