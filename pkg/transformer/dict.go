package transformer

import "io"

type Dict struct {
	Name   string
	Reader io.Reader
}

type Entry struct {
	Word  string
	Code  string
	Order int
}
