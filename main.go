package main

import (
	"context"
	"github.com/emilhauk/chitchat/config"
	"github.com/emilhauk/chitchat/internal/controller"
	"github.com/emilhauk/chitchat/internal/database"
	"github.com/emilhauk/chitchat/internal/manager"
	internalMiddleware "github.com/emilhauk/chitchat/internal/middleware"
	"github.com/emilhauk/chitchat/internal/server"
)

var (
	log            = config.Logger
	err            error
	dbStore        database.DBStore
	userManager    manager.User
	sessionManager manager.Session
	channelManager manager.Channel
	messageManager manager.Message
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewConnectionPool(config.Database)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer db.Close()

	dbStore := database.NewDBStore(db)
	userManager = manager.NewUserManager(dbStore.Users, dbStore.Credentials)
	sessionManager = manager.NewSessionManager(dbStore.Sessions)
	channelManager = manager.NewChannelManager(dbStore.Channels)
	messageManager = manager.NewMessageManager(dbStore.Messages)

	// TODO This stinks. Should provide better wrapper for controllers
	controller.ProvideManagers(userManager, sessionManager, channelManager, messageManager)

	authMiddleware := internalMiddleware.NewAuthMiddleware(userManager, sessionManager)
	router := server.NewRouter(authMiddleware)

	server.Start(ctx, router)
}
