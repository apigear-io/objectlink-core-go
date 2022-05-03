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

func (w *MockDataWriter) WriteData(data []byte) error {
	msg, err := w.converter.FromData(data)
	if err != nil {
		return err
	}
	w.Messages = append(w.Messages, msg)
	return nil
}
