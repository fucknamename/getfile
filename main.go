package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"getfile/utils"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	key_token            = "fuck9527code"
	TIME_LAYOUT_FORMAT_2 = "20060102150405"
)

var (
	regex     = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	sensitive = []string{"etc", "root", "var", "www", "home", "tmp", "proc", "sys"}
)

func init() {
	setToken()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/gq/9527", handleToken)
	mux.HandleFunc("/getfile/*", handleFile)

	srv := &http.Server{
		Addr:    ":1688",
		Handler: mux,
	}

	fmt.Println("file server run at 1688 port")
	fmt.Println("telegram: @echoty")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

// 读取token
func handleToken(w http.ResponseWriter, r *http.Request) {
	token, ok := utils.GetCache(key_token)
	if ok {
		w.Write([]byte(token.(string)))
	} else {
		setToken()
		w.Write([]byte("try again"))
	}
}

// 打包下载文件
func handleFile(w http.ResponseWriter, r *http.Request) {
	var (
		/*
			http://xx.xx.xx.xx:xx/getfile/ff_b/release/aadddddd
		*/
		path  = strings.TrimPrefix(r.URL.Path, "/")
		parts = strings.Split(path, "/")
	)

	if len(parts) != 4 || !strings.HasPrefix(parts[0], "getfile") {
		http.NotFound(w, r)
		return
	}

	if parts[1] == "" || parts[2] == "" || parts[1] == "code" ||
		strings.Contains(parts[1], ".") || strings.Contains(parts[2], ".") {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("illegal request"))
		return
	}

	if !regex.MatchString(parts[1]) || !regex.MatchString(parts[2]) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error path 001"))
		return
	}

	for _, v := range sensitive {
		if parts[1] == v || parts[2] == v {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("error path 002"))
			return
		}
	}

	token, ok := utils.GetCache(key_token)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("inner error"))
		return
	}

	if parts[3] != token.(string) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("token not correct"))
		return
	}

	file, err := utils.ArchiveZip(fmt.Sprintf("%s/%s", parts[1], parts[2]))
	if err != nil {
		w.Write([]byte("archive zip faild"))
	} else if file == "" {
		w.Write([]byte("no compiled files"))
	} else {
		// 直接下载文件
		http.ServeFile(w, r, file)
	}
}

func setToken() {
	dt := time.Now().Format(TIME_LAYOUT_FORMAT_2)
	val := getMD5Hash(dt + "ty")
	utils.SetCache(key_token, val)
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
