package got

import (
	"context"
	"sync"
)

type handlerRegistration struct {
	gen  HandlerGenerator
	opts []MuxOption
}

type funcRegistration struct {
	gen  HandlerFuncGenerator
	opts []MuxOption
}

var (
	once           sync.Once
	hRegistrations []*handlerRegistration
	fRegistrations []*funcRegistration
)

func RegisterHandler(hg HandlerGenerator, opts ...MuxOption) {
	once.Do(func() { hRegistrations = []*handlerRegistration{} })
	hRegistrations = append(hRegistrations, &handlerRegistration{hg, opts})
}

func RegisterHandlerFunc(hg HandlerFuncGenerator, opts ...MuxOption) {
	once.Do(func() { fRegistrations = []*funcRegistration{} })
	fRegistrations = append(fRegistrations, &funcRegistration{hg, opts})
}

func Start(ctx context.Context, token string, opts ...Option) error {
	once.Do(func() { hRegistrations = []*handlerRegistration{} })
	once.Do(func() { fRegistrations = []*funcRegistration{} })

	got, err := New(ctx, token, opts...)
	if err != nil {
		return err
	}

	for _, hr := range hRegistrations {
		if err := got.RegisterHandler(hr.gen, hr.opts...); err != nil {
			return err
		}
	}

	for _, fr := range fRegistrations {
		if err := got.RegisterHandlerFunc(fr.gen, fr.opts...); err != nil {
			return err
		}
	}

	return got.Start(ctx)
}
