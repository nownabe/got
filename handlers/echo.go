package handlers

import (
	"context"
	"regexp"

	"github.com/nownabe/got"
)

func Echo() got.HandlerGenerator {
	return func(g *got.Got) (got.Handler, error) {
		re := regexp.MustCompile(`\Aecho`)
		f := func(ctx context.Context, msg *got.Message) {
			msg.Reply(msg.Text())
		}

		mux := g.NewMux()
		mux.AddHandlerFunc(f, got.HandleOnlyMention(), got.HandleWithRegexp(re))

		return mux, nil
	}
}
