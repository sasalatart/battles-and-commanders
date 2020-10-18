package logger

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Interface defines behaviour for "info" and "error" logging
type Interface interface {
	Info(text string)
	Error(err error)
}

type logger struct {
	infoWriter io.Writer
	errWriter  io.Writer
}

// New returns a logger.Interface instance by using two different io.Writer values: one for info
// logs and another one for error logs
func New(infoWriter, errWriter io.Writer) Interface {
	return logger{infoWriter: infoWriter, errWriter: errWriter}
}

// NewDiscard returns a logger.Interface compatible logger which discards every commanded write
func NewDiscard() Interface {
	return New(ioutil.Discard, ioutil.Discard)
}

// Info writes a text message into the service's info io.Writer
func (s logger) Info(text string) {
	s.infoWriter.Write([]byte(text))
}

// Error writes an error message into the service's error io.Writer
func (s logger) Error(err error) {
	s.errWriter.Write([]byte(fmt.Sprintln(err)))
}
