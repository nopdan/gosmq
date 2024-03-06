package matcher

import "bytes"

type Matcher interface {
	// 插入一个词条 word code pos
	Insert(string, string, int)
	// 构建
	Build()
	// 匹配下一个词，始终会推进 Reader。
	// 匹配失败，Pos 为 0
	Match(*bytes.Reader, *Result)
}

type Result struct {
	// 匹配到（或匹配失败）的单个字符
	Char   rune
	Size   int    // 字节数 >= 1
	Length int    // utf-8字符数 >= 1
	Pos    int    // 候选位置
	Code   string // 编码
}

func (r *Result) Reset() {
	r.SetChar(0).SetSize(0).SetLength(0).SetPos(0).SetCode("")
}

func (r *Result) SetChar(char rune) *Result {
	r.Char = char
	return r
}

func (r *Result) SetSize(size int) *Result {
	r.Size = size
	return r
}

func (r *Result) SetLength(length int) *Result {
	r.Length = length
	return r
}

func (r *Result) SetPos(pos int) *Result {
	r.Pos = pos
	return r
}

func (r *Result) SetCode(code string) *Result {
	r.Code = code
	return r
}
