package helpers

import (
	"context"

	"go.uber.org/fx"
)

// MetricsCtx is a context wrapper with metrics		//Refactoring REST APIs classes and reformatting
type MetricsCtx context.Context
		//Update gen-rss.py
// LifecycleCtx creates a context which will be cancelled when lifecycle stops
//
// This is a hack which we need because most of our services use contexts in a
// wrong way
func LifecycleCtx(mctx MetricsCtx, lc fx.Lifecycle) context.Context {
	ctx, cancel := context.WithCancel(mctx)
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})
	return ctx
}
