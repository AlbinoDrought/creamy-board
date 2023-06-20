package log

import stdlog "log"

func Warnf(format string, v ...any) {
	stdlog.Printf(format, v...)
}

func Infof(format string, v ...any) {
	stdlog.Printf(format, v...)
}
