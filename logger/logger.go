package logger

type Logger interface {
	Debugf(string, ...interface{})
	Debug(...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
	Fatalf(string, ...interface{})
	Fatal(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Infof(string, ...interface{})
	Info(...interface{})
}

