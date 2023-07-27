package logger

type ILogger interface {
	Info(msg string)
	Warning(msg string)
}
