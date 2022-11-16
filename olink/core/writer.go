package core

type DataWriter interface {
	WriteData(data []byte) error
}
