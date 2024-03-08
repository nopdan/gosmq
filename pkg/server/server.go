package server

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/nopdan/gosmq/pkg/util"
)

//go:embed dist
var dist embed.FS

// 上传的文件列表
var files [][]byte = make([][]byte, 0)

func Serve(port string, silent bool) {
	mux := http.NewServeMux()
	dist, _ := fs.Sub(dist, "dist")
	mux.Handle("GET /", http.FileServer(http.FS(dist)))
	mux.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		type Result struct {
			Text []string
			Dict []string
		}
		// text := util.WalkDirWithSuffix("./text/", ".txt")
		// dict := util.WalkDirWithSuffix("./dict/", ".txt")
		text := util.WalkDirWithSuffix(`D:\Code\go\gosmq\build\text`, ".txt")
		dict := util.WalkDirWithSuffix(`D:\Code\go\gosmq\build\dict`, ".txt")
		res := Result{Text: text, Dict: dict}
		b, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("GET /list error: %v\n", err)
			return
		}
		_, _ = w.Write(b)
	})
	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		// 获取上传的文件
		err := r.ParseMultipartForm(1024 * 1024 * 1024) // 最大 1GB
		if err != nil {
			log.Printf("POST /upload err: %v\n", err)
			return
		}
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Printf("POST /upload err: %v\n", err)
			return
		}
		defer file.Close()
		log.Printf("POST /upload:%v", handler.Header)

		data, err := io.ReadAll(file)
		if err != nil {
			log.Printf("POST /upload err: %v\n", err)
			return
		}
		index := len(files)
		files = append(files, data)
		_, _ = w.Write([]byte("{\"index\": " + strconv.Itoa(index) + "}"))
	})
	mux.HandleFunc("POST /race", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		text := r.FormValue("text")
		dict := r.FormValue("dict")
		fmt.Printf("    text: %v\n    dict: %v\n", text, dict)

		rd := r.Body
		b, err := io.ReadAll(rd)
		if err != nil {
			return
		}
		fmt.Printf("POST /race body: %v\n", string(b))
		_, _ = w.Write([]byte("{\"hello\": \"world\"}"))
	})

	var wg sync.WaitGroup
	wg.Add(1)
	port = ":" + port
	url := "http://localhost" + port
	go func() {
		fmt.Println("Listen and serve: ", url)
		err := http.ListenAndServe(port, mux)
		if err != nil {
			fmt.Println(err)
			panic("...Serve failed.")
		}
		wg.Done()
	}()
	if !silent {
		openBrowser(url)
	}
	wg.Wait()
}

func setHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
