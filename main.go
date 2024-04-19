package main

import (
	"context"
	"flag"
	"fmt"
	"getfile/handle"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	_ = handle.SetToken()
}

func main() {
	tplpath := flag.String("t", "", "待检查的模板文件夹路径")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/tp/check", func(w http.ResponseWriter, r *http.Request) {
		handle.HandleCheck(w, r, *tplpath)
	})
	mux.HandleFunc("/gq/9527", handle.HandleToken)
	mux.HandleFunc("/getfile/*", handle.HandleFile)

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
