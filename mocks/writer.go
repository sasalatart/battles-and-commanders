package mocks

// Writer mocks an io.Writer interface, storing each one of the messages written in string form
type Writer struct {
	Writes []string
}

func (w *Writer) Write(bytes []byte) (int, error) {
	w.Writes = append(w.Writes, string(bytes))
	return 0, nil
}
