package iface

type Ilogger interface {
	Info(string)
	Error(string)
	Panic(string)
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Json(any) string
}
