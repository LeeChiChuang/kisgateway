package logx

import "log"

type LoggerWriter struct {
	logger *log.Logger
}

func NewLoggerWriter(l *log.Logger) LoggerWriter {
	return LoggerWriter{logger: l}
}

func (l LoggerWriter) Write(data []byte) (n int, err error) {
	l.logger.Print(string(data))
	return len(data), nil
}

func (l LoggerWriter) Close() error {
	return nil
}
