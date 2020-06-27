package mocks

// Logger mocks an io.Writer interface, storing each one of the messages written in string form
type Logger struct {
	Logs []string
}

func (l *Logger) Write(bytes []byte) (int, error) {
	l.Logs = append(l.Logs, string(bytes))
	return 0, nil
}
