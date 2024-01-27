package utils

import (
	"context"
	"sync"

	"time"
)

type Handler[T any, V any] struct {
	ctx    context.Context
	cancel context.CancelFunc
	Ich    chan *T
	Och    chan *V
	Wg     *sync.WaitGroup
}

func newHandler[T any, V any](ctx context.Context, cancel context.CancelFunc) *Handler[T, V] {
	return &Handler[T, V]{
		ctx:    ctx,
		cancel: cancel,
		Ich:    make(chan *T),
		Och:    make(chan *V),
		Wg:     &sync.WaitGroup{},
	}
}

func NewHandlerWithTimeout[T any, V any](timeout time.Duration) *Handler[T, V] {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return newHandler[T, V](ctx, cancel)
}

func NewHandlerFromContext[T any, V any](ctx context.Context) *Handler[T, V] {
	ctx, cancel := context.WithCancel(ctx)
	return newHandler[T, V](ctx, cancel)
}

func NewHandler[T any, V any]() *Handler[T, V] {
	ctx, cancel := context.WithCancel(context.Background())
	return newHandler[T, V](ctx, cancel)
}

func (h *Handler[T, V]) Ctx() context.Context {
	return h.ctx
}

func (h *Handler[T, V]) Cancel() {
	h.cancel()
}

func (h *Handler[T, V]) SetInputChannel(ich chan *T) {
	h.Ich = ich
}

func (h *Handler[T, V]) SetOutputChannel(och chan *V) {
	h.Och = och
}
