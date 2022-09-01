package transformer

import "io"

type Dict struct {
	SavePath string
	Name     string
	Reader   io.Reader
}

type Entry struct {
	Word  string
	Code  string
	Order int
}
