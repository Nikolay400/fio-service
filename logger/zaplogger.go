package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

type ZapLooger struct {
	l *zap.Logger
}

func NewZapLogger() (*ZapLooger, error) {
	zaplogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	zaplogger.Sync()
	return &ZapLooger{zaplogger}, nil
}

func (zl *ZapLooger) Info(s string) {
	zl.l.Info(s)
}

func (zl *ZapLooger) Error(s string) {
	zl.l.Error(s)
}

func (zl *ZapLooger) Panic(s string) {
	zl.l.Panic(s)
}

func (zl *ZapLooger) Json(obj any) string {
	var jsonStr string
	jsonByte, err := json.Marshal(obj)
	if err != nil {
		jsonStr = err.Error()
	} else {
		jsonStr = string(jsonByte)
	}
	return jsonStr
}

func (zl *ZapLooger) Infof(template string, args ...interface{}) {
	zl.l.Sugar().Infof(template, args...)
}

func (zl *ZapLooger) Errorf(template string, args ...interface{}) {
	zl.l.Sugar().Errorf(template, args...)
}

func (zl *ZapLooger) Panicf(template string, args ...interface{}) {
	zl.l.Sugar().Panicf(template, args...)
}
