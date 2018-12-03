package got

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/nlopes/slack"
	"github.com/nownabe/got/store"
	"github.com/nownabe/got/stores/memory"
	"go.uber.org/zap"
)

type Got struct {
	id        string
	logger    *zap.SugaredLogger
	mentionRE *regexp.Regexp
	name      string
	rootMux   *Mux
	rtm       *slack.RTM
	slack     Slack
	store     store.Provider
	wg        *sync.WaitGroup
}

func New(ctx context.Context, token string, opts ...Option) (*Got, error) {
	return NewWithSlack(ctx, slack.New(token), opts...)
}

func NewWithSlack(ctx context.Context, s Slack, opts ...Option) (got *Got, err error) {
	got = &Got{
		slack: s,
		wg:    &sync.WaitGroup{},
	}
	got.rootMux = got.NewMux()

	l, err := zap.NewProduction()
	if err != nil {
		return
	}
	got.logger = l.Sugar()
	got.store = memory.New()

	for _, o := range opts {
		o.apply(got)
	}

	if err := got.getBotInfo(ctx); err != nil {
		got.logger.Errorw("failed to get bot info", "error", err)
		return got, err
	}

	got.mentionRE, err = regexp.Compile(fmt.Sprintf(`\A\s*<@%s>\s*`, got.id))
	if err != nil {
		return
	}

	return
}

func (got *Got) getBotInfo(ctx context.Context) error {
	resp, err := got.slack.AuthTestContext(ctx)
	if err != nil {
		return err
	}
	got.id = resp.UserID
	got.name = resp.User
	return nil
}

func (got *Got) ID() string {
	return got.id
}

func (got *Got) Logger() *zap.SugaredLogger {
	return got.logger
}

func (got *Got) Name() string {
	return got.name
}

func (got *Got) AddHandler(h Handler, opts ...MuxOption) {
	got.rootMux.AddHandler(h, opts...)
}

func (got *Got) AddHandlerFunc(f HandlerFunc, opts ...MuxOption) {
	got.AddHandler(f, opts...)
}

func (got *Got) RegisterHandler(hgen HandlerGenerator, opts ...MuxOption) error {
	h, err := hgen(got)
	if err != nil {
		got.logger.Errorw("failed to generate handler", "error", err)
		return err
	}
	got.AddHandler(h, opts...)
	return nil
}

func (got *Got) RegisterHandlerFunc(hgen HandlerFuncGenerator, opts ...MuxOption) error {
	f, err := hgen(got)
	if err != nil {
		got.logger.Errorw("failed to generate handler func", "error", err)
		return err
	}
	got.AddHandlerFunc(f, opts...)
	return nil
}

func (got *Got) Start(ctx context.Context) error {
	got.rtm = got.slack.NewRTM()
	go got.rtm.ManageConnection()

	got.logger.Infow("started")
	for {
		select {

		case <-ctx.Done():
			got.wg.Wait()
			if err := got.rtm.Disconnect(); err != nil {
				got.logger.Warnw("failed to disconnect", "error", err)
			}
			return ctx.Err()

		case msg := <-got.rtm.IncomingEvents:
			got.logger.Infow("received event", "event", msg)

			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				got.handle(ctx, ev)

			case *slack.LatencyReport:
				got.logger.Infow("latency report", "latency", ev.Value)

			case *slack.RTMError:
				got.logger.Errorw("rtm error", "error", ev.Error())

			case *slack.InvalidAuthEvent:
				got.logger.Errorw("invalid authentication")
				return errors.New("invalid authentication")

			default:
				// Ignore
			}
		}
	}
}

func (got *Got) Store(namespace string) store.Store {
	return got.store.Store(namespace)
}

func (got *Got) handle(ctx context.Context, e *slack.MessageEvent) {
	msg := got.newMessage(e)
	got.rootMux.Handle(ctx, msg)
}
