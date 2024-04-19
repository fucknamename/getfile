package handle

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"getfile/utils"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	key_token          = "fuck9527code"
	TIME_LAYOUT_FORMAT = "20060102150405"
)

var (
	regex     = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	sensitive = []string{"etc", "root", "var", "www", "home", "tmp", "proc", "sys"}
)

// 校验模板文件
func HandleCheck(w http.ResponseWriter, r *http.Request, tplpath string) {
	if tplpath == "" {
		w.Write([]byte("not config tpl directory, use -t 'directory'"))
		return
	}

	if dangerTpl := RunTplSafeTasks(tplpath); dangerTpl != "" {
		w.Write([]byte(fmt.Sprintf("模板文件 %s 存在安全隐患，请检查", dangerTpl)))
		return
	} else {
		w.Write([]byte("1"))
	}
}

// 读取token
func HandleToken(w http.ResponseWriter, r *http.Request) {
	token, ok := utils.GetCache(key_token)
	if ok {
		w.Write([]byte(token.(string)))
	} else {
		t := SetToken()
		w.Write([]byte(t))
	}
}

// 打包下载文件
func HandleFile(w http.ResponseWriter, r *http.Request) {
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

func SetToken() string {
	dt := time.Now().Format(TIME_LAYOUT_FORMAT)
	val := getMD5Hash(dt + "ty")
	utils.SetCache(key_token, val)

	return val
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
