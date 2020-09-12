package logger

import (
	"fmt"
	"io"
)

// Service is a wrapper that provides logging functionality, separating info logs from error logs
type Service struct {
	infoWriter io.Writer
	errWriter  io.Writer
}

// New creates a new instance of logger.Service by using two different io.Writer values: one for
// info logs and another one for error logs
func New(infoWriter, errWriter io.Writer) Service {
	return Service{infoWriter: infoWriter, errWriter: errWriter}
}

// Info writes a text message into the service's info io.Writer
func (s Service) Info(text string) {
	s.infoWriter.Write([]byte(text))
}

// Error writes an error message into the service's error io.Writer
func (s Service) Error(err error) {
	s.errWriter.Write([]byte(fmt.Sprintln(err)))
}
