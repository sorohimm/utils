// Package log contains a set of functions and structs for creating and working with zap logger
package log

import (
	"context"
	stdl "log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	UserId    = "userId"
	RequestId = "reqId"
	Source    = "source"
)

var (
	global     *Zap
	globalOnce = sync.Once{}
)

func GetGlobal() *Zap {
	globalOnce.Do(func() {
		l, err := NewZap("debug", "json")
		if err != nil {
			stdl.Fatal(err)
		}
		global = l
	})

	return global
}

// NewZap creates a new instance of zap logger with the provided level and encoding type.
// lvl: level of the logger (debug, info, warn, error, etc)
// encType: encoding type of the logger (json, console, etc)
// Returns: new instance of zap logger and error if any
func NewZap(lvl, encType string) (*Zap, error) {
	var (
		err    error
		logger *zap.Logger
	)

	config := zap.NewProductionConfig()
	config.Encoding = encType
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "c"
	config.EncoderConfig.StacktraceKey = "s"

	if err = config.Level.UnmarshalText([]byte(lvl)); err != nil && lvl != "" {
		return nil, err
	}

	if logger, err = config.Build(); err != nil {
		return nil, err
	}

	return &Zap{logger}, nil
}

type Zap struct {
	*zap.Logger
}

func (o *Zap) WithServerOptions(name, version, built string) *Zap {
	return &Zap{
		Logger: o.WithOptions(zap.AddStacktrace(zapcore.DPanicLevel)).
			With(zap.String("v", version), zap.String("built", built), zap.String("app", name)),
	}
}

func (o *Zap) WithSource(source string) *Zap {
	return &Zap{
		Logger: o.Logger.With(zap.String(Source, source)),
	}
}

func (o *Zap) WithRequestId(id string) *Zap {
	return &Zap{
		Logger: o.Logger.With(zap.String(RequestId, id)),
	}
}

// ctxKey is a struct that is used as the key for storing logger in the context
type ctxKey struct{} // or exported to use outside the package

// CtxWithLogger adds the provided logger to the context and returns the new context
func CtxWithLogger(ctx context.Context, logger *Zap) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext retrieves the logger from the context, returns a Nop logger if not found
func FromContext(ctx context.Context) *Zap {
	if ctxLogger, ok := ctx.Value(ctxKey{}).(*Zap); ok {
		return ctxLogger
	}

	return &Zap{zap.NewNop()}
}
