package handlers

import (
	"context"
	"regexp"

	"github.com/nownabe/got"
)

func Ping() got.HandlerGenerator {
	return func(g *got.Got) (got.Handler, error) {
		re := regexp.MustCompile(`\Aping`)
		f := func(ctx context.Context, msg *got.Message) {
			msg.Reply("pong")
		}

		mux := g.NewMux()
		mux.AddHandlerFunc(f, got.HandleOnlyMention(), got.HandleWithRegexp(re))

		return mux, nil
	}
}
