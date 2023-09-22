package matcher

import (
	"fmt"
	"slices"
	"time"
)

type trie struct {
	root   *trieNode
	values []value

	tails   []tail
	useTail bool // 是否压缩 tail

	count  uint32 // 插入词的数量
	stable bool   // 是否按照码表的顺序
}

type trieNode struct {
	ch map[rune]*trieNode

	valueIdx int32
	tailIdx  int32
	pass     uint32 // 经过节点的次数
}

type value struct {
	code  string
	pos   int
	order uint32 // 插入节点的顺序
}

type tail struct {
	runes    []rune
	valueIdx int32
}

func NewTrie(stable bool, useTail bool) *trie {
	t := new(trie)
	t.root = newTrieNode()
	t.values = make([]value, 0, 1e4)
	t.stable = stable
	if useTail {
		t.useTail = useTail
		t.tails = make([]tail, 0, 1000)
	}
	return t
}

func newTrieNode() *trieNode {
	tn := new(trieNode)
	tn.valueIdx = -1
	tn.tailIdx = -1
	return tn
}

func (t *trie) Insert(word, code string, pos int) {
	node := t.root
	for _, v := range word {
		if node.ch == nil {
			node.ch = make(map[rune]*trieNode)
			node.ch[v] = newTrieNode()
		} else if node.ch[v] == nil {
			node.ch[v] = newTrieNode()
		}
		node.pass++
		node = node.ch[v]
	}
	t.count++
	// 新词
	if node.valueIdx == -1 {
		node.valueIdx = int32(len(t.values))
		t.values = append(t.values, value{code, pos, t.count})
		return
	}
	// 已经存在的词
	// 取排在前面的
	if t.stable {
		return
	}
	// 取码长较短的
	value := &t.values[node.valueIdx]
	if len(value.code) > len(code) {
		value.code = code
		value.pos = pos
		value.order = t.count
	}
}

func (t *trie) Build() {
	if t.useTail {
		start := time.Now()
		node := t.root
		node.build(&t.tails)
		fmt.Printf("构建 tail 耗时: %dms\n", time.Since(start).Milliseconds())
	}
}

func (node *trieNode) build(tails *[]tail) {
	if node.ch == nil {
		return
	}
	if node.pass == 1 {
		node.mergeTail(tails)
		return
	}
	for _, ch := range node.ch {
		ch.build(tails)
	}
}

// 取唯一的孩子节点
func getUniqueNode(node map[rune]*trieNode) (rune, *trieNode) {
	if len(node) != 1 {
		panic("children node not unique")
	}
	for rn, ch := range node {
		return rn, ch
	}
	return 0, nil
}

// 合并 tail 节点
func (head *trieNode) mergeTail(tails *[]tail) {
	rn, node := getUniqueNode(head.ch)
	// 单字 tail
	// AB ABC   B->C(tail)
	if node.ch == nil {
		return
	}
	// 多字 tail
	// AB ABCD  B->CD(tail)
	runes := []rune{rn}
	for node.ch != nil {
		rn, node = getUniqueNode(node.ch)
		runes = append(runes, rn)
	}
	head.ch = nil
	head.tailIdx = int32(len(*tails))
	*tails = append(*tails, tail{runes, node.valueIdx})
}

// 前缀树最长匹配
func (t *trie) Match(text []rune) (int, string, int) {

	node := t.root
	wordLen := 0
	res := new(value)

	match := func(p int, _tail tail) {
		if p+len(_tail.runes) > len(text) {
			return
		}
		if slices.Equal(_tail.runes, text[p:p+len(_tail.runes)]) {
			val := &t.values[_tail.valueIdx]
			// 跳过码表顺序在后面的词
			if t.stable && res.order != 0 && val.order > res.order {
				return
			}
			wordLen = p + len(_tail.runes)
			res = val
		}
	}

	for p := 0; p < len(text); {
		node = node.ch[text[p]]
		p++
		if node == nil {
			break
		}
		if node.valueIdx != -1 {
			val := &t.values[node.valueIdx]
			// 跳过码表顺序在后面的词
			if t.stable && res.order != 0 && val.order > res.order {
			} else {
				wordLen = p
				res = val
			}
		}

		// 匹配 tail
		if t.useTail && node.tailIdx != -1 {
			_tail := t.tails[node.tailIdx]
			match(p, _tail)
			break
		}
	}
	return wordLen, res.code, res.pos
}
