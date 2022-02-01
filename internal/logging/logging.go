package logging

import (
	logRotator "github.com/KaiserWerk/go-log-rotator"
)

type Logger struct {
	rotator *logRotator.Rotator
}

func New(path string) (*Logger, error) {
	rotator, err := logRotator.New(path, "loggy.log", 10<<20, 0644, 0, false)
	if err != nil {
		return nil, err
	}

	return &Logger{
		rotator: rotator,
	}, nil
}

func (l *Logger) Write(b []byte) (int, error) {
	return l.rotator.Write(b)
}

func (l *Logger) Close() error {
	return l.rotator.Close()
}
