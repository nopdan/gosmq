package web2

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"

	_ "embed"

	"github.com/cxcn/gosmq/pkg/smq"
)

//go:embed index.html
var index string

func Run() {

	type names struct {
		DictNames []string
		TextNames []string
	}

	theNames := new(names)
	// 读取dict/目录中的所有文件和子目录
	files, err := os.ReadDir(`dict/`)
	if err != nil {
		fmt.Println("找不到 dict 文件夹", err)
		return
	}
	fmt.Println("检测到以下赛码表：")
	for _, file := range files {
		theNames.DictNames = append(theNames.DictNames, file.Name())
		fmt.Println("  ", file.Name())
	}
	fmt.Println()
	// 读取text/目录中的所有文件和子目录
	files, err = os.ReadDir(`text/`)
	if err != nil {
		fmt.Println("找不到 text 文件夹", err)
		return
	}
	fmt.Println("检测到以下文本：")
	for _, file := range files {
		theNames.TextNames = append(theNames.TextNames, file.Name())
		fmt.Println("  ", file.Name())
	}

	funcMap := template.FuncMap{"getName": smq.GetFileName}
	t := template.New("index.html").Funcs(funcMap)
	_, err = t.Parse(index)
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr: "localhost:5666",
	}
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// html, _ := ioutil.ReadFile("index.html")
		// rw.Write(html)
		t.Execute(rw, theNames)
	})
	http.HandleFunc("/result", func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		options := getOptions(r.PostForm)
		if len(options) == 0 {
			return
		}
		tn := smq.GetFileName(r.PostForm.Get("fpt"))
		h := NewHTML(tn)
		s := smq.NewFromPath(tn, "./text/"+r.PostForm.Get("fpt"))
		for _, v := range options {
			s.Add(v)
		}
		res := s.Run()
		for _, v := range res {
			h.AddResult(v)
		}
		h.OutputHTML(rw)

		// fmt.Fprintln(rw, r.PostForm)
	})
	fmt.Println("\nhttp://localhost:5666/\n ")
	server.ListenAndServe()
}

func getOptions(v url.Values) []*smq.Dict {
	var ret []*smq.Dict
	if v.Get("fpd") != "" {
		tmp := new(smq.Dict)
		tmp.LoadFromPath("./dict/" + v.Get("fpd"))

		tmp.PushStart, _ = strconv.Atoi(v.Get("ding"))
		tmp.SelectKeys = v.Get("csk")
		if v.Get("iss") == "true" {
			tmp.Single = true
		}
		tmp.PressSpaceBy = "both"
		ret = append(ret, tmp)
	}
	if v.Get("fpd1") != "" {
		tmp := new(smq.Dict)
		tmp.LoadFromPath("./dict/" + v.Get("fpd1"))

		tmp.PushStart, _ = strconv.Atoi(v.Get("ding1"))
		tmp.SelectKeys = v.Get("csk1")
		if v.Get("iss1") == "true" {
			tmp.Single = true
		}
		tmp.PressSpaceBy = "both"
		ret = append(ret, tmp)
	}
	return ret
}
