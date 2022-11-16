package core

type MockDataWriter struct {
	Messages  []Message
	converter *MessageConverter
}

func NewMockDataWriter() *MockDataWriter {
	return &MockDataWriter{
		converter: &MessageConverter{Format: FormatJson},
	}
}

func (w *MockDataWriter) Write(data []byte) (int, error) {
	msg, err := w.converter.FromData(data)
	if err != nil {
		return 0, err
	}
	w.Messages = append(w.Messages, msg)
	return len(data), nil
}

func (w *MockDataWriter) Close() error {
	return nil
}
