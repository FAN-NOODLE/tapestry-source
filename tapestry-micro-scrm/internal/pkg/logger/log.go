package logger

import (
	"github.com/rs/zerolog"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

var Log *Logger

// Logger 放入模块的base中 从base中取
type Logger struct {
}

func NewLogger() *Logger {
	Log = &Logger{}
	return Log
}

// Trace 必须调用 Msg 方法才能输出日志
func (l Logger) Trace(cc ctx.Context, options ...log.Option) *zerolog.Event {
	return log.Trace(cc, options...)
}

// Debug 必须调用 Msg 方法才能输出日志
func (l Logger) Debug(cc ctx.Context, options ...log.Option) *zerolog.Event {
	return log.Debug(cc, options...)
}

// Info 必须调用 Msg 方法才能输出日志
func (l Logger) Info(cc ctx.Context, options ...log.Option) *zerolog.Event {
	return log.Info(cc, options...)
}

// Warn 必须调用 Msg 方法才能输出日志
func (l Logger) Warn(cc ctx.Context, options ...log.Option) *zerolog.Event {
	return log.Warn(cc, options...)
}

// Error 必须调用 Msg 方法才能输出日志
func (l Logger) Error(cc ctx.Context, options ...log.Option) *zerolog.Event {
	return log.Error(cc, options...)
}
