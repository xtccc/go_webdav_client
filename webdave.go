package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"errors"

	"github.com/jedib0t/go-pretty/table"
	"github.com/schollz/progressbar/v3"
	"github.com/studio-b12/gowebdav"
)

func download(c *gowebdav.Client, webdavFilePath, localFilePath string) {
	// 判断一下，如果 localFilePath 这个是目录 则自己将 webdavFilePath 中的文件名 拼接到 localFilePath 后面

	lf, err := os.Stat(localFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

		} else {
			fmt.Println(err)
		}
	} else {
		if lf.IsDir() {
			filename := filepath.Base(webdavFilePath)
			localFilePath = filepath.Join(localFilePath, filename)
		}
	}

	//获取一下远程文件大小
	fileinfo, err := c.Stat(webdavFilePath)
	if err != nil {
		fmt.Println(err)
	}

	bar := progressbar.DefaultBytes(fileinfo.Size(), "downloading")
	reader, _ := c.ReadStream(webdavFilePath)
	// 使用 TeeReader 包装文件读取器
	bar_reader := io.TeeReader(reader, bar)

	file, _ := os.Create(localFilePath)
	defer file.Close()

	io.Copy(file, bar_reader)

}

func delete(c *gowebdav.Client, webdavFilePath string) {
	err := c.Remove(webdavFilePath)
	if err != nil {
		fmt.Println(err)
	}
}

func mkdir(c *gowebdav.Client, webdavFilePath string) {
	err := c.MkdirAll(webdavFilePath, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
func listfile(c *gowebdav.Client, webdavFilePath string) {
	fmt.Println("list ", webdavFilePath)

	files, err := c.ReadDir(webdavFilePath)
	if err != nil {
		if gowebdav.IsErrCode(err, 405) {
			file, err := c.Stat(webdavFilePath)
			if err != nil {
				fmt.Println(err)
			}
			files = append(files, file)

		} else {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	sort.Slice(files,

		func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"F/D", "Name", "Size", "Modified"})
	for _, file := range files {
		Type := '-'
		size := "-"
		time := file.ModTime().Format("2006-01-02 15:04:05")
		if file.IsDir() {
			Type = 'd'

		} else {
			size = size_format(file.Size())
		}
		t.AppendRow([]interface{}{string(Type), file.Name(), size, time})
	}
	t.Render()

}
func upload(c *gowebdav.Client, localFilePath, webdavFilePath string) {
	// 判断一下，如果 webdavFilePath 这个是目录 则自己将localFilePath中的文件名 拼接到 webdavFilePath后面

	webf, err := c.Stat(webdavFilePath)
	if err != nil {
		if gowebdav.IsErrCode(err, 404) {
			//远程是文件并且不存在 继续走即可
		} else {
			fmt.Println(err)
		}

	} else {
		//fmt.Println("webf.IsDir", webf.IsDir())
		if webf.IsDir() {
			filename := filepath.Base(localFilePath)
			webdavFilePath = filepath.Join(webdavFilePath, filename)
		}

	}

	fileinfo, err := os.Stat(localFilePath)
	if err != nil {
		fmt.Println(err)
	}

	bar := progressbar.DefaultBytes(fileinfo.Size(), "uploading")

	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// 使用 TeeReader 包装文件读取器
	reader := io.TeeReader(file, bar)

	err = c.WriteStream(webdavFilePath, reader, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func size_format(bytes int64) string {
	// 输入数字
	//输出大小 string
	const (
		kib = 1024
		mib = kib * 1024
		gib = mib * 1024
	)

	switch {
	case bytes >= gib:
		return fmt.Sprintf("%.2fG", float64(bytes)/gib)
	case bytes >= mib:
		return fmt.Sprintf("%.2fM", float64(bytes)/mib)
	case bytes >= kib:
		return fmt.Sprintf("%.2fK", float64(bytes)/kib)
	default:
		return fmt.Sprintf("%dB", bytes)
	}

}

func init_url(url string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	configPath := filepath.Join(homeDir, ".config", "webdav.conf")
	// 创建文件所在的目录（如果不存在）
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	f, err := os.Create(configPath)
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString(url)
}

func get_url() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
	}
	configPath := filepath.Join(homeDir, ".config", "webdav.conf")
	f_byte, err := os.ReadFile(configPath)
	if err != nil {
		return "", err

	}
	url := string(f_byte)
	return url, err
}
