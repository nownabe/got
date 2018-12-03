package got

import (
	"context"
	"regexp"
	"sync"
)

type Mux struct {
	handlers []*muxHandler
	wg       *sync.WaitGroup
}

func (got *Got) NewMux() *Mux {
	return &Mux{
		handlers: []*muxHandler{},
		wg:       got.wg,
	}
}

func (m *Mux) AddHandler(h Handler, opts ...MuxOption) {
	mh := &muxHandler{
		handler:    h,
		allMessage: true,
		rawRegexps: []*regexp.Regexp{},
		regexps:    []*regexp.Regexp{},
	}

	for _, o := range opts {
		o.apply(mh)
	}

	m.handlers = append(m.handlers, mh)
}

func (m *Mux) AddHandlerFunc(f HandlerFunc, opts ...MuxOption) {
	m.AddHandler(f, opts...)
}

func (m *Mux) Handle(ctx context.Context, msg *Message) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		for _, h := range m.handlers {
			m.wg.Add(1)
			go func(h *muxHandler) {
				defer m.wg.Done()
				h.Handle(ctx, msg)
			}(h)
		}
	}()
}

type MuxOption interface {
	apply(*muxHandler)
}

func HandleAllMessages() MuxOption {
	return muxOptionFunc(func(h *muxHandler) {
		h.allMessage = true
	})
}

func HandleOnlyMention() MuxOption {
	return muxOptionFunc(func(h *muxHandler) {
		h.allMessage = false
	})
}

func HandleWithRegexp(re *regexp.Regexp) MuxOption {
	return muxOptionFunc(func(h *muxHandler) {
		h.regexps = append(h.regexps, re)
	})
}

func HandleWithRegexpRaw(re *regexp.Regexp) MuxOption {
	return muxOptionFunc(func(h *muxHandler) {
		h.rawRegexps = append(h.rawRegexps, re)
	})
}

type muxOptionFunc func(h *muxHandler)

func (f muxOptionFunc) apply(h *muxHandler) {
	f(h)
}

type muxHandler struct {
	handler Handler

	// options
	allMessage bool
	rawRegexps []*regexp.Regexp
	regexps    []*regexp.Regexp
}

func (h *muxHandler) Handle(ctx context.Context, msg *Message) {
	if !h.allMessage && !msg.IsMention() {
		return
	}

	if !h.matchRE(msg) {
		return
	}

	h.handler.Handle(ctx, msg)
}

func (h *muxHandler) matchRE(msg *Message) bool {
	if len(h.regexps) == 0 && len(h.rawRegexps) == 0 {
		return true
	}

	for _, re := range h.rawRegexps {
		if re.MatchString(msg.FullText()) {
			return true
		}
	}

	for _, re := range h.regexps {
		if re.MatchString(msg.Text()) {
			return true
		}
	}

	return false
}
