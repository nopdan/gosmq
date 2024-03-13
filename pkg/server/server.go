package server

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/nopdan/gosmq/pkg/util"
)

//go:embed dist
var dist embed.FS

// 上传的文件列表
var files [][]byte = make([][]byte, 0)

var textList = make([]string, 0)

func Serve(port int, silent bool, prefix string) {
	mux := http.NewServeMux()
	dist, _ := fs.Sub(dist, "dist")
	mux.Handle("GET /", http.FileServer(http.FS(dist)))
	mux.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		logger.Info("GET /list")
		res := struct {
			Text    []string `json:"text"`
			Dict    []string `json:"dict"`
			TextDir string   `json:"textDir"`
			DictDir string   `json:"dictDir"`
		}{
			TextDir: filepath.Join(prefix, "text"),
			DictDir: filepath.Join(prefix, "dict"),
		}
		res.Text = util.WalkDirWithSuffix(res.TextDir, ".txt")
		textList = res.Text
		res.Dict = util.WalkDirWithSuffix(res.DictDir, ".txt")
		json.NewEncoder(w).Encode(res)
	})
	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		logger.Info("POST /upload")
		// 获取上传的文件
		err := r.ParseMultipartForm(1024 * 1024 * 1024) // 最大 1GB
		if err != nil {
			logger.With("error", err).Error("POST /upload")
			return
		}
		file, handler, err := r.FormFile("file")
		if err != nil {
			logger.With("error", err).Error("POST /upload")
			return
		}
		defer file.Close()
		logger.Info("POST /upload", "header", handler.Header)

		data, err := io.ReadAll(file)
		if err != nil {
			logger.With("error", err).Error("POST /upload")
			return
		}
		files = append(files, data)
	})
	mux.HandleFunc("GET /file_index", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		logger.Info("GET /file_index")
		index := len(files) - 1
		_, _ = w.Write([]byte("{\"index\": " + strconv.Itoa(index) + "}"))
	})
	mux.HandleFunc("POST /race", func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w)
		logger.Info("POST /race")
		data := r.FormValue("data")

		d := &Data{}
		err := json.Unmarshal([]byte(data), d)
		if err != nil {
			logger.With("error", err).Error("POST /race")
			return
		}
		res := d.Race()
		_, _ = w.Write(res)
	})

	var wg sync.WaitGroup
	wg.Add(1)

	url := fmt.Sprintf("http://localhost:%d", port)
	go func() {
		fmt.Println("Listen and serve: ", url)
		addr := ":" + strconv.Itoa(port)
		err := http.ListenAndServe(addr, mux)
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
