package got

import (
	"github.com/nownabe/got/store"
	"go.uber.org/zap"
)

type Option interface {
	apply(*Got)
}

type optionFunc func(*Got)

func (f optionFunc) apply(got *Got) {
	f(got)
}

func WithLogger(l *zap.SugaredLogger) Option {
	return optionFunc(func(got *Got) {
		got.logger = l
	})
}

func WithStore(s store.Provider) Option {
	return optionFunc(func(got *Got) {
		got.store = s
	})
}
