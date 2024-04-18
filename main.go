package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"getfile/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
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
	_ = setToken()
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

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefulExitWeb(srv)
}

// 优雅退出
func gracefulExitWeb(server *http.Server) {
	quit := make(chan os.Signal, 4)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit

	fmt.Println("got a signal\ngetfile server stoped", sig)

	now := time.Now()
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(cxt); err != nil {
		fmt.Println("err", err)
	}

	// 看看实际退出所耗费的时间
	fmt.Println("------exited--------", time.Since(now))
	os.Exit(0)
}

// 读取token
func handleToken(w http.ResponseWriter, r *http.Request) {
	token, ok := utils.GetCache(key_token)
	if ok {
		w.Write([]byte(token.(string)))
	} else {
		t := setToken()
		w.Write([]byte(t))
	}
}

// 打包下载文件
func handleFile(w http.ResponseWriter, r *http.Request) {
	var (
		/*
			http://xx.xx.xx.xx:xx/getfile/package/token字符串
		*/
		dirlen = 10 // 文件夹名长度限制
		path   = strings.TrimPrefix(r.URL.Path, "/")
		parts  = strings.Split(path, "/")
	)

	if len(parts) != 3 || !strings.HasPrefix(parts[0], "getfile") {
		http.NotFound(w, r)
		return
	}

	if parts[1] == "" || parts[2] == "" || parts[1] == "code" || strings.Contains(parts[1], ".") {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("illegal request"))
		return
	}

	if len(parts[1]) > dirlen {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("too loooooooooog")) //target project dir name too long
		return
	}

	if !regex.MatchString(parts[1]) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error path 001"))
		return
	}

	for _, v := range sensitive {
		if parts[1] == v {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("error path 002"))
			return
		}
	}

	token, ok := utils.GetCache(key_token)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("johnny knows everything"))
		return
	}

	if parts[2] != token.(string) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("token not correct"))
		return
	}

	file, err := utils.ArchiveZip(parts[1])
	if err != nil {
		w.Write([]byte("archive zip faild"))
	} else if file == "" {
		w.Write([]byte("no compiled files"))
	} else {
		// 直接下载文件
		http.ServeFile(w, r, file)
	}
}

func setToken() string {
	dt := time.Now().Format(TIME_LAYOUT_FORMAT_2)
	val := getMD5Hash(dt + "ty")
	utils.SetCache(key_token, val)

	return val
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
