package core

type MockDataWriter struct {
	Messages []Message
	conv     *MessageConverter
}

func NewMockDataWriter() *MockDataWriter {
	return &MockDataWriter{
		conv: &MessageConverter{Format: FormatJson},
	}
}

func (w *MockDataWriter) Write(data []byte) (int, error) {
	msg, err := w.conv.FromData(data)
	if err != nil {
		return 0, err
	}
	w.Messages = append(w.Messages, msg)
	return len(data), nil
}

func (w *MockDataWriter) Close() error {
	return nil
}
