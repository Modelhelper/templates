package zap_logger

import (
	"{{ .Name | kebab }}/service"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type zaplogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger(config *zap.Config) (service.LogService, error) {
	logger, err := config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "port.LogService.NewZapLogger")
	}
	return &zaplogger{
		logger.Sugar(),
	}, nil
}

func (l *zaplogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}
func (l *zaplogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}
func (l *zaplogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}
func (l *zaplogger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}
func (l *zaplogger) Warning(args ...interface{}) {
	l.sugar.Warn(args...)
}
func (l *zaplogger) Warningf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}
func (l *zaplogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}
func (l *zaplogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}
func (l *zaplogger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}
func (l *zaplogger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}
