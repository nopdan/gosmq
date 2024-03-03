package serve

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/nopdan/gosmq/pkg/smq"
)

//go:embed dist
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
	Path   string `json:"path"`
	Single bool   `json:"single"`
	Space  string `json:"space"`
	Stable bool   `json:"stable"`
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
		Stable:       opt.Stable,
		PressSpaceBy: opt.Space,
	}
	dict.Load("dict/" + opt.Path)
	fmt.Println("载入码表：", dict.Name)
	return dict
}

var smqRes []*smq.Result

func GetResultJson(src []byte) []byte {
	var opts = parseOptions(src)
	s := &smq.Text{}
	if opts.Text.Flag {
		if opts.Text.Name == "" {
			opts.Text.Name = "赛文"
		}
		s.LoadString(opts.Text.Name, opts.Text.Plain)
		fmt.Println("载入文本：", opts.Text.Name)
	} else {
		s.Load("text/" + opts.Text.Path)
		fmt.Println("载入文本：", s.Name)
	}
	dicts := make([]*smq.Dict, 0)
	for _, v := range opts.Dicts {
		// opt := toSmqDict(v)
		dicts = append(dicts, toSmqDict(v))
		// s.Add(opt)
	}
	smqRes = s.Race(dicts, false)
	result, _ := json.Marshal(smqRes)
	return result
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	setHeader(&w)
	defer r.Body.Close()
	if r.Method == "POST" {
		start := time.Now()
		body, _ := io.ReadAll(r.Body)
		// fmt.Println("    post body: ", string(body))
		rjson := GetResultJson(body)
		w.Write(rjson)
		// fmt.Println("    returned json: ", string(rjson))
		fmt.Printf("比赛结束，耗时：%v\n\n", time.Since(start))
	}
}

// 递归遍历文件夹
func getFiles(dirname string, pre string) []string {
	ret := make([]string, 0)

	fileInfos, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			//继续遍历fi这个目录
			tmp := getFiles(dirname+"/"+fi.Name(), fi.Name()+"/")
			ret = append(ret, tmp...)
		} else {
			ret = append(ret, pre+fi.Name())
		}
	}
	return ret
}

func setHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func Serve(port string, silent bool) {
	dist, _ := fs.Sub(dist, "dist")
	http.Handle("/", http.FileServer(http.FS(dist)))
	http.HandleFunc("/race", RaceHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/file_index", IndexHandler)
	http.HandleFunc("/api", PostHandler)
	http.HandleFunc("/texts", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		b, _ := json.Marshal(getFiles(`text/`, ""))
		w.Write(b)
	})
	http.HandleFunc("/dicts", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		b, _ := json.Marshal(getFiles(`dict/`, ""))
		w.Write(b)
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		if len(smqRes) == 0 {
			return
		}
		h := NewHTML()
		for _, v := range smqRes {
			h.AddResult(v)
		}
		h.OutputHTML(w)
	})

	url := "http://localhost:" + port
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println("Listen and serve: ", url)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Println(err)
			panic("...Server failed.")
		}
		wg.Done()
	}()
	if !silent {
		// openBrowser(url)
	}
	wg.Wait()
}

func openBrowser(url string) {
	var name string
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
}

func RaceHandler(w http.ResponseWriter, r *http.Request) {
	setHeader(&w)
	text := r.FormValue("text")
	dict := r.FormValue("dict")
	fmt.Printf("    text: %v\n    dict: %v\n", text, dict)
	w.Write([]byte("{\"hello\": \"world\"}"))
}

var fileList = make([][]byte, 0, 1)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	setHeader(&w)
	r.ParseMultipartForm(1024 * 1024 * 1024) // 最大 1GB
	fmt.Println(r.Form)

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	defer file.Close()
	log.Printf("POST:%v", handler.Header)

	mbData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	fileList = append(fileList, mbData)
	// fmt.Println(string(mbData))
	index := len(fileList) - 1
	w.Write([]byte("{'index': " + strconv.Itoa(index) + "}"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	setHeader(&w)
	index := len(fileList)
	w.Write([]byte("{\"index\": " + strconv.Itoa(index) + "}"))
}
