package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/mholt/archiver/v3"
)

// 解压rar文件
func UnarchiveRar(rarfile, dir string) error {
	unarchiver := archiver.NewRar()     // 创建一个 RAR 类型的解压缩器
	unarchiver.OverwriteExisting = true // 覆盖
	err := unarchiver.Unarchive(rarfile, dir)
	if err != nil {
		return err
	}

	// 删除压缩包文件
	go func(r, d string) {
		if err = os.Remove(r); err != nil {
			fmt.Println(err)
		}
	}(rarfile, dir)

	return nil
}

// 压缩zip文件
func ArchiveZip(path string) (string, error) {
	var (
		err  error
		file = "ff.zip"
	)

	path = "/home/" + strings.TrimPrefix(path, "/") // 限制只能是home目录下

	ach := archiver.NewZip()     // 创建一个 ZIP 类型的解压缩器
	ach.OverwriteExisting = true // 如果文件存在也要重新打包
	err = ach.Archive([]string{path}, fmt.Sprintf("%s/%s", path, file))

	return path + "/" + file, err
}

// 解压rar文件
func UnarchiveZip(rarfile, dir string) error {
	unarchiver := archiver.NewZip()     // 创建一个 ZIP 类型的解压缩器
	unarchiver.OverwriteExisting = true // 覆盖
	err := unarchiver.Unarchive(rarfile, dir)
	if err != nil {
		return err
	}

	// 删除压缩包文件
	go func(r, d string) {
		if err = os.Remove(r); err != nil {
			fmt.Println(err)
		}
	}(rarfile, dir)

	return nil
}

// 解压tar.gz文件
func UnarchiveTarGz(rarfile, dir string) error {
	unarchiver := archiver.NewTarGz()   // 创建一个 tar.gz 类型的解压缩器
	unarchiver.OverwriteExisting = true // 覆盖
	err := unarchiver.Unarchive(rarfile, dir)
	if err != nil {
		return err
	}

	// 删除压缩包文件
	go func(r, d string) {
		if err = os.Remove(r); err != nil {
			fmt.Println(err)
		}
	}(rarfile, dir)

	return nil
}
