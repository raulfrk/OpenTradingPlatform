package utils

import (
	"context"
	"fmt"
	"sync"

	"time"
)

type Handler[T any] struct {
	ctx                    context.Context
	cancel                 context.CancelFunc
	Ch                     chan *T
	RoundRobinChannels     []chan *T
	Wg                     *sync.WaitGroup
	Lock                   sync.RWMutex
	FunctionalityAttatched bool
	currentAgent           int
}

func newHandler[T any](ctx context.Context, cancel context.CancelFunc) *Handler[T] {
	return &Handler[T]{
		ctx:    ctx,
		cancel: cancel,
		Ch:     make(chan *T),
		Wg:     &sync.WaitGroup{},
	}
}

func newRoundRobinHandler[T any](ctx context.Context, cancel context.CancelFunc, maxAgents int) *Handler[T] {
	var agents []chan *T
	for i := 0; i < maxAgents; i++ {
		agents = append(agents, make(chan *T))
	}
	return &Handler[T]{
		ctx:                ctx,
		cancel:             cancel,
		Ch:                 make(chan *T),
		RoundRobinChannels: agents,
		Wg:                 &sync.WaitGroup{},
	}
}

func (h *Handler[T]) GetNextAgent() (chan *T, error) {
	h.Lock.Lock()
	if h.RoundRobinChannels == nil {
		return nil, fmt.Errorf("handler is not a round robin handler")
	}
	if h.currentAgent == len(h.RoundRobinChannels) || len(h.RoundRobinChannels) == 1 {
		h.currentAgent = 1
	} else {
		h.currentAgent++
	}
	h.Lock.Unlock()
	return h.RoundRobinChannels[h.currentAgent-1], nil

}

func NewHandlerWithTimeout[T any, V any](timeout time.Duration) *Handler[T] {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return newHandler[T](ctx, cancel)
}

func NewHandlerFromContext[T any, V any](ctx context.Context) *Handler[T] {
	ctx, cancel := context.WithCancel(ctx)
	return newHandler[T](ctx, cancel)
}

func NewHandler[T any]() *Handler[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return newHandler[T](ctx, cancel)
}

func NewRoundRobinHandler[T any](maxAgents int) *Handler[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return newRoundRobinHandler[T](ctx, cancel, maxAgents)
}

func (h *Handler[T]) Ctx() context.Context {
	return h.ctx
}

func (h *Handler[T]) Cancel() {
	h.cancel()
}

func (h *Handler[T]) SetChannel(ich chan *T) {
	h.Ch = ich
}
