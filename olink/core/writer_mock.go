package core

type MockWriter struct {
	Messages []Message
}

func NewMockWriter() *MockWriter {
	return &MockWriter{
		Messages: make([]Message, 0),
	}
}

func (w *MockWriter) WriteMessage(m Message) error {
	w.Messages = append(w.Messages, m)
	return nil
}
