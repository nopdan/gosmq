package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func read(fp string) Trie {

	var dict = Constructor()
	dict.addPunct()
	conf := readConf(fp)

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("read cost time = ", cost)
	}()

	if conf.isConf {
		fmt.Println("检测到普通码表", conf)
	} else {
		fmt.Println("检测到赛码表")
		dict.readSMB(fp)
		return dict
	}

	f, err := os.Open(fp)
	errHandler(err)
	defer f.Close()

	freq := make(map[string]int)
	buff := bufio.NewReader(f)
	// 生成字典
	for {
		b, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		wc := strings.Split(string(b), "\t")
		if len(wc) != 2 {
			continue
		}
		word, code := wc[0], wc[1]
		freq[code] += 1
		var key string
		if freq[code] != 1 {
			key = code + strconv.Itoa(freq[code])
		} else if len(code) < conf.as {
			key = code + "_"
		} else {
			key = code
		}
		dict.Insert(word, key)
	}
	return dict
}

func (t Trie) readSMB(fp string) {
	f, err := ioutil.ReadFile(fp)
	errHandler(err)
	fs := string(f)
	fs = strings.Replace(fs, "\r", "", -1)
	words := strings.Split(fs, "\n")
	for _, v := range words {
		wc := strings.Split(v, "\t")
		if len(wc) == 2 {
			t.Insert(wc[0], wc[1])
		}
	}
}

func (t Trie) addPunct() {
	// 符号
	punct := map[string]string{
		"·": "`",
		"【": "[",
		"】": "]",
		"；": ";",
		"‘": "\"",
		"’": "\"",
		"，": ",",
		"。": ".",
		"、": "/",
		// 用 ~ 表示换挡
		"——": "~-",
		"：":  "~;",
		"“":  "~\"",
		"”":  "~\"",
		"《":  "~,",
		"》":  "~.",
		"？":  "~/",
	}
	for k, v := range punct {
		t.Insert(k, v)
	}
}
