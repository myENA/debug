package sync

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
}
