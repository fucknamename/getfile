package handle

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// golang 危险标记关键词
	danger_mark = []string{"[%", "%]", "package ", "import", "func"}
)

// 模板文件安全性校验任务
func RunTplSafeTasks(tpldir string) string {
	var (
		count int      // 模板文件数量
		tpls  []string // 模板html文件集
		wg    sync.WaitGroup
	)

	err := filepath.Walk(tpldir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			tpls = append(tpls, path)
		}
		return nil
	})

	if count = len(tpls); count == 0 || err != nil {
		return ""
	}
	if count == 1 { // 直接判断
		errtags := 0
		filepath := tpls[0]
		checkcount := len(danger_mark)
		content, _ := ReadFromFile(filepath)
		htmlstr := string(content)

		if htmlstr == "" {
			return ""
		}

		for i := 0; i < checkcount; i++ {
			if strings.Contains(danger_mark[i], "exec") || strings.Contains(danger_mark[i], "Exec") {
				errtags++
			}
			if strings.Contains(htmlstr, danger_mark[i]) {
				errtags++
			}
		}

		if errtags >= checkcount { // 全部命中
			return filepath
		} else {
			return ""
		}
	}

	// 消息通道
	ch := make(chan string, count)

	// 启动校验任务
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			t := validateTplPage(tpls[idx])
			if t != "" {
				ch <- t
			}
		}(i)
	}

	// 等待所有任务完成或者校验失败
	go func() {
		wg.Wait()
		close(ch) // 关闭通道，确保所有任务完成后关闭通道
	}()

	// 检查是否有任务出错
	if tpl := <-ch; tpl != "" {
		return tpl
	}

	return ""
}

// 校验模板文件
func validateTplPage(filepath string) string {
	// 读取文件
	errtags := 0
	count := len(danger_mark)
	content, _ := ReadFromFile(filepath)
	htmlstr := string(content)

	if htmlstr == "" {
		return ""
	}

	for i := 0; i < count; i++ {
		if strings.Contains(danger_mark[i], "exec") || strings.Contains(danger_mark[i], "Exec") {
			errtags++
		}
		if strings.Contains(htmlstr, danger_mark[i]) {
			errtags++
		}
	}

	if errtags >= len(danger_mark) { // 全部命中
		return filepath
	} else {
		return ""
	}
}

// ReadFromFile 从文件中读取内容
func ReadFromFile(filename string) ([]byte, error) {
	// 读取不同的文件，这里不需要加锁
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err //文件不存在
	}
	return os.ReadFile(filename)
}
