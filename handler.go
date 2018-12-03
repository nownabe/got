package got

import (
	"context"
)

type Handler interface {
	Handle(context.Context, *Message)
}

type HandlerFunc func(context.Context, *Message)

func (f HandlerFunc) Handle(ctx context.Context, msg *Message) {
	f(ctx, msg)
}

type HandlerGenerator func(*Got) (Handler, error)

type HandlerFuncGenerator func(*Got) (HandlerFunc, error)
