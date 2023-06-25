package handlers

import (
	"context"
	"github.com/thuongnn/clst-mgt-api/models"
	"log"
	"sync"
)

type Handler func(message *models.EventMessage) error

type MessageHandler struct {
	ctx      context.Context
	mutex    sync.RWMutex
	handlers map[models.EventType]Handler
}

func NewMessageHandler(ctx context.Context) *MessageHandler {
	return &MessageHandler{
		ctx:      ctx,
		handlers: make(map[models.EventType]Handler),
	}
}

func (mh *MessageHandler) RegisterHandler(messageType models.EventType, handler Handler) {
	mh.mutex.Lock()
	defer mh.mutex.Unlock()

	mh.handlers[messageType] = handler
}

func (mh *MessageHandler) HandleMessage(message *models.EventMessage) error {
	handler, ok := mh.handlers[message.Type]
	if !ok {
		log.Printf("No handler found for message type: %s\n", message.Type)
		return nil
	}

	return handler(message)
}
