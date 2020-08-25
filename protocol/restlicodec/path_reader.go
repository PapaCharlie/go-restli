package restlicodec

type PathReader interface {
	Reader
	ReadPathSegment() string
}

func NewPathReader() PathReader {
	panic("implement me!")
}
