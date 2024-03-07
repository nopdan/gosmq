package matcher

import (
	"bytes"
	"io"
)

type Trie struct {
	root    *trieNode
	count   int  // 插入词的数量
	ordered bool // 是否按照码表的顺序
}

type value struct {
	code  string
	pos   int
	order int
}

type trieNode struct {
	ch    map[rune]*trieNode
	value *value
}

func NewTrie(ordered bool) *Trie {
	t := new(Trie)
	t.root = new(trieNode)
	t.root.ch = make(map[rune]*trieNode, 8192)
	t.ordered = ordered
	return t
}

func (t *Trie) Insert(word, code string, pos int) {
	node := t.root
	for _, v := range word {
		if node.ch == nil {
			node.ch = make(map[rune]*trieNode)
			node.ch[v] = &trieNode{}
		} else if _, ok := node.ch[v]; !ok {
			node.ch[v] = &trieNode{}
		}
		node = node.ch[v]
	}
	t.count++
	if node.value == nil {
		node.value = &value{code, pos, t.count}
	} else if !t.ordered {
		// 贪心，已经存在的词，取码长较短的
		if len(node.value.code) > len(code) {
			node.value.code = code
			node.value.pos = pos
			node.value.order = t.count
		}
	}
}

func (t *Trie) Build() {}

func (t *Trie) Match(brd *bytes.Reader, res *Result) {
	res.Reset()
	node := t.root
	res.Length = 1 // 至少匹配一个字
	var Char rune
	var CharSize int
	var Size, Length int
	var order int
	for {
		char, size, err := brd.ReadRune()
		if err != nil {
			break
		}
		if Char == 0 {
			Char = char
			CharSize = size
		}
		Size += size
		Length++

		node = node.ch[char]
		if node == nil {
			break
		}
		if node.value != nil {
			if !t.ordered || node.value.order > order {
				order = node.value.order
				res.Size = Size
				res.Length = Length
				res.Code = node.value.code
				res.Pos = node.value.pos
			}
		}
	}
	if res.Length == 1 {
		res.Char = Char
		brd.Seek(int64(CharSize-Size), io.SeekCurrent)
	} else {
		brd.Seek(int64(res.Size-Size), io.SeekCurrent)
	}
}
