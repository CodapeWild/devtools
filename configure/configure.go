package configure

type Configure interface {
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
}
