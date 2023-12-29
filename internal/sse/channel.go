package sse

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/emilhauk/chitchat/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"sync"
)

type ChatService interface {
	IsMemberOfChannel(channelUUID, userUUID string) (bool, error)
	GetChannelList(user model.User) ([]model.Channel, error)
}

type Event struct {
	ID              string
	Type            string
	Message         model.Message
	Channel         model.Channel
	CurrentUserUUID string
}

func NewEvent(t string, channel model.Channel, message model.Message, currentUserUUID string) Event {
	id := uuid.NewString()

	return Event{
		ID:              id,
		Type:            t,
		Channel:         channel,
		Message:         message,
		CurrentUserUUID: currentUserUUID,
	}
}

type Broker struct {
	channelConsumers     map[chan Event]string
	channelListConsumers map[chan Event]string
	logger               zerolog.Logger
	chatService          ChatService
	mtx                  *sync.Mutex
}

func NewBroker(logger zerolog.Logger, chatService ChatService) *Broker {
	return &Broker{
		channelConsumers:     make(map[chan Event]string),
		channelListConsumers: make(map[chan Event]string),
		chatService:          chatService,
		mtx:                  new(sync.Mutex),
		logger:               logger,
	}
}

func (b *Broker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), app.BrokerContextKey, b)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (b *Broker) SubscribeToChannel(channelUUID string) chan Event {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	c := make(chan Event)
	b.channelConsumers[c] = channelUUID

	b.logger.Debug().Msgf("client connected to channel %s", channelUUID)
	return c
}

func (b *Broker) SubscribeToChannelListUpdates(userUUID string) chan Event {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	c := make(chan Event)
	b.channelListConsumers[c] = userUUID

	b.logger.Debug().Msgf("client connected to channelList %s", userUUID)
	return c
}

func (b *Broker) Unsubscribe(c chan Event) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	id := b.channelConsumers[c]
	if id != "" {
		close(c)
		delete(b.channelConsumers, c)
		b.logger.Debug().Msgf("Client %s killed, %d remaining", id, len(b.channelConsumers))
	}

	id = b.channelListConsumers[c]
	if id != "" {
		close(c)
		delete(b.channelListConsumers, c)
		b.logger.Debug().Msgf("Client %s killed, %d remaining", id, len(b.channelListConsumers))
	}
}

func (b *Broker) Publish(e Event) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	go func() {
		// TODO this truly sucks. It will perform a lot of checks which could've been avoided if channel just included a list of its members.
		for s, _ := range b.channelListConsumers {
			s <- e
		}
	}()

	pubMsg := 0
	for s, channelUUID := range b.channelConsumers {
		if channelUUID == e.Channel.UUID {
			s <- e
			pubMsg++
		}
	}

	b.logger.Debug().Msgf("published message to %d subscribers", pubMsg)
}

func (b *Broker) Close() {
	for k := range b.channelConsumers {
		close(k)
		delete(b.channelConsumers, k)
	}
}

func (b *Broker) ServeHTTPForChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported", http.StatusInternalServerError)
		return
	}

	user := app.GetUserFromContextOrPanic(r.Context())

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create new client channel for stream events
	c := b.SubscribeToChannel(chi.URLParam(r, "channelUUID"))
	defer b.Unsubscribe(c)

	for {
		select {
		case msg := <-c:
			if msg.Message.Sender.UUID == user.UUID {
				continue
			}
			// TODO We generate message fom template for each recipient here. This seems inefficient.
			buf := bytes.Buffer{}
			err := templates.Templates.ExecuteTemplate(&buf, "message", msg.Message)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to execute template")
				continue
			}
			_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Type, strings.ReplaceAll(buf.String(), "\n", ""))
			f.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func (b *Broker) ServeHTTPForChannelList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported", http.StatusInternalServerError)
		return
	}

	user := app.GetUserFromContextOrPanic(r.Context())

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create new client channel for stream events
	c := b.SubscribeToChannelListUpdates(user.UUID)
	defer b.Unsubscribe(c)

	for {
		select {
		case msg := <-c:
			_, err := b.chatService.IsMemberOfChannel(msg.Channel.UUID, user.UUID)
			if err != nil {
				if !errors.Is(err, app.ErrMemberNotFound) {
					b.logger.Error().Err(err).Msgf("Failed sending channel=(%s) list update to userUUID=(%s)", msg.Channel.UUID, user.UUID)
				}
				continue
			}
			channelList, err := b.chatService.GetChannelList(user)
			if err != nil {
				b.logger.Error().Err(err).Msgf("Failed to get channel list for userUUID=(%s)", user.UUID)
				continue
			}
			buf := bytes.Buffer{}
			err = templates.Templates.ExecuteTemplate(&buf, "channel-list", channelList)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to execute template")
				continue
			}
			_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", "channelList", strings.ReplaceAll(buf.String(), "\n", ""))
			f.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func PublishUsingBrokerInContext(ctx context.Context, event Event) error {
	if broker, ok := ctx.Value(app.BrokerContextKey).(*Broker); ok {
		broker.Publish(event)
		return nil
	}
	return errors.New("no message broker found in context")
}
