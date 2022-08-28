package transformer

import "io"

type Dict struct {
	SavePath   string
	Name       string
	Reader     io.Reader
	PushStart  int
	SelectKeys string
	Single     bool
}
