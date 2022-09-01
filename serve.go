package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/cxcn/gosmq/pkg/smq"
)

//go:embed web/dist
var dist embed.FS

type Options struct {
	Text  optText   `json:"text"`
	Dicts []optDict `json:"dicts"`
	Web   string    `json:"web"`
}

type optText struct {
	Name  string `json:"name"`
	Plain string `json:"plain"`
	Path  string `json:"path"`
	Flag  bool   `json:"flag"` // 手动输入赛文
}

type optDict struct {
	Path       string `json:"path"`
	Format     string `json:"format"`
	Single     bool   `json:"single"`
	SelectKeys string `json:"selectkeys"`
	PushStart  int    `json:"pushstart"`
	Alg        string `json:"alg"`
}

func parseOptions(src []byte) Options {
	var opt Options
	if e := json.Unmarshal(src, &opt); e != nil {
		fmt.Println("...Error in Unmarshal: ", e.Error())
	}
	return opt
}

func toSmqDict(opt optDict) *smq.Dict {
	dict := &smq.Dict{
		Single:       opt.Single,
		Format:       opt.Format,
		SelectKeys:   opt.SelectKeys,
		PushStart:    opt.PushStart,
		Algorithm:    opt.Alg,
		PressSpaceBy: "both",
		OutputDetail: false,
		OutputDict:   false,
	}
	dict.LoadFromPath("./dict/" + opt.Path)
	return dict
}

var smqRes []*smq.Result
var textName string

func GetResultJson(src []byte) []byte {
	var opts = parseOptions(src)
	var s smq.Smq
	if opts.Text.Flag {
		if opts.Text.Name == "" {
			opts.Text.Name = "赛文"
		}
		s = smq.NewFromString(opts.Text.Name, opts.Text.Plain)
	} else {
		opts.Text.Name = smq.GetFileName(opts.Text.Path)
		s = smq.NewFromPath(opts.Text.Name, "./text/"+opts.Text.Path)
	}
	textName = opts.Text.Name
	for _, v := range opts.Dicts {
		opt := toSmqDict(v)
		s.Add(opt)
	}
	smqRes = s.Run()
	result, _ := smq.ResToJson(smqRes)
	return result
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	setHeader(&w)
	defer r.Body.Close()
	if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		fmt.Println("    post body: ", string(body))
		rjson := GetResultJson(body)
		w.Write(rjson)
		fmt.Println("    returned json: ", string(rjson))
	}
}

// 读取目录中的所有文件
func getFiles(path string) []string {
	ret := make([]string, 0, 1)
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("找不到：", path, err)
		return ret
	}
	for _, file := range files {
		ret = append(ret, file.Name())
	}
	return ret
}

func setHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func serve() {
	fsys, _ := fs.Sub(dist, "web/dist")
	http.Handle("/", http.FileServer(http.FS(fsys)))
	http.HandleFunc("/api", PostHandler)
	http.HandleFunc("/texts", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		b, _ := json.Marshal(getFiles(`text/`))
		w.Write(b)
	})
	http.HandleFunc("/dicts", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		b, _ := json.Marshal(getFiles(`dict/`))
		w.Write(b)
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		if len(smqRes) == 0 {
			return
		}
		h := NewHTML(textName)
		for _, v := range smqRes {
			h.AddResult(v)
		}
		h.OutputHTML(w)
	})

	var wg sync.WaitGroup
	wg.Add(1)
	port := ":7172"
	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			panic("...Server failed.")
		}
		wg.Done()
	}()
	doSomething(port)
	wg.Wait()
}

func doSomething(port string) {
	var name string
	url := "http://localhost" + port
	switch runtime.GOOS {
	case "windows":
		name = "explorer"
	case "linux":
		name = "xdg-open"
	default:
		name = "open"
	}
	cmd := exec.Command(name, url)
	cmd.Start()
	fmt.Println("Listen and serve: ", url)
}
