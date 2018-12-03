package handlers

import (
	"context"
	"regexp"

	"github.com/nownabe/got"
	"github.com/nownabe/got/store"
)

func KVS(namespace string) got.HandlerGenerator {
	return func(g *got.Got) (got.Handler, error) {
		mux := g.NewMux()

		getRE := regexp.MustCompile(`\Aget\s+([^\s]+)\z`)
		getF := func(ctx context.Context, msg *got.Message) {
			sub := getRE.FindStringSubmatch(msg.Text())
			key := sub[1]

			var reply string
			val, err := g.Store(namespace).GetString(ctx, key)
			switch err {
			case store.ErrStoreNotFound:
				reply = key + " is not found"
			case store.ErrStoreNotStringValue:
				reply = "invalid value"
			case nil:
				reply = val
			default:
				reply = "error"
			}
			msg.Reply(reply)
		}
		mux.AddHandlerFunc(getF, got.HandleOnlyMention(), got.HandleWithRegexp(getRE))

		setRE := regexp.MustCompile(`\Aset\s+([^\s]+)\s+(.+)\z`)
		setF := func(ctx context.Context, msg *got.Message) {
			sub := setRE.FindStringSubmatch(msg.Text())
			key := sub[1]
			val := sub[2]

			var reply string
			err := g.Store(namespace).Set(ctx, key, val)
			if err == nil {
				reply = "ok"
			} else {
				reply = "failed to set " + key
			}
			msg.Reply(reply)
		}
		mux.AddHandlerFunc(setF, got.HandleOnlyMention(), got.HandleWithRegexp(setRE))

		return mux, nil
	}
}
