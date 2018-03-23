package sync

import (
	stdlog "log"
	"os"
)

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
}

var log Logger

func init() {
	log = stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
}

func SetLogger(l Logger) {
	log = l
}
