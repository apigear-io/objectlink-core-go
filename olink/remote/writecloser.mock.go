package remote

type MockWriteCloser struct {
	Messages     [][]byte
	WriteHandler func(p []byte) (n int, err error)
	Closed       bool
	CloseHandler func() error
}

func NewMockWriteCloser() *MockWriteCloser {
	return &MockWriteCloser{
		Messages: make([][]byte, 0),
	}
}

func (m *MockWriteCloser) Write(p []byte) (n int, err error) {
	m.Messages = append(m.Messages, p)
	if m.WriteHandler != nil {
		return m.WriteHandler(p)
	}
	return len(p), nil
}

func (m *MockWriteCloser) Close() error {
	m.Closed = true
	if m.CloseHandler != nil {
		return m.CloseHandler()
	}
	return nil
}
