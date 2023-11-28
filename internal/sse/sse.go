package sse

import (
	"context"
	"errors"
	"fmt"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
	"sync"
)

type Event struct {
	ID             string
	Type           string
	Message        string
	OriginUserUUID string
	ChannelUUID    string
}

func NewEvent(t string, channelUUID string, originUserUUID, message string) Event {
	id := uuid.NewString()

	return Event{
		ID:             id,
		Type:           t,
		Message:        message,
		OriginUserUUID: originUserUUID,
		ChannelUUID:    channelUUID,
	}
}

type Broker struct {
	consumers map[chan Event]string
	logger    zerolog.Logger
	mtx       *sync.Mutex
}

func NewBroker(logger zerolog.Logger) *Broker {
	return &Broker{
		consumers: make(map[chan Event]string),
		mtx:       new(sync.Mutex),
		logger:    logger,
	}
}

func (b *Broker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), app.BrokerContextKey, b)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (b *Broker) Subscribe(channelUUID string) chan Event {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	c := make(chan Event)
	b.consumers[c] = channelUUID

	b.logger.Debug().Msgf("client connected to channel %s", channelUUID)
	return c
}

func (b *Broker) Unsubscribe(c chan Event) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	id := b.consumers[c]
	close(c)
	delete(b.consumers, c)
	b.logger.Debug().Msgf("Client %s killed, %d remaining", id, len(b.consumers))
}

func (b *Broker) Publish(e Event) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	pubMsg := 0
	for s, channelUUID := range b.consumers {
		// Push to specific channel
		if channelUUID == e.ChannelUUID {
			s <- e
			pubMsg++
		}
	}

	b.logger.Debug().Msgf("published message to %d subscribers", pubMsg)
}

func (b *Broker) Close() {
	for k := range b.consumers {
		close(k)
		delete(b.consumers, k)
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	c := b.Subscribe(chi.URLParam(r, "channelUUID"))
	defer b.Unsubscribe(c)

	for {
		select {
		case msg := <-c:
			if msg.OriginUserUUID == user.UUID {
				continue
			}
			_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Type, strings.ReplaceAll(msg.Message, "\n", ""))
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
