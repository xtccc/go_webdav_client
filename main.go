package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/studio-b12/gowebdav"
)

func webdav_init() *gowebdav.Client {
	root, err := get_url()
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("请使用init 命令初始化配置url")
		os.Exit(1)
	}
	user := "user"
	password := "password"

	c := gowebdav.NewClient(root, user, password)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
	}

	c.SetTransport(tr)
	err = c.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return c

}
func main() {

	var rootCmd = &cobra.Command{
		//是命令的名称
		Use: "go_webdave",
		//表示当前命令的完整描述。
		Long: "go_webdave 是一个上传下载工具",

		//Run属性是一个函数，当执行命令时会调用此函数。
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "上传文件",
		Long:  "上传文件 需要两个参数：本地文件路径 和 webdav文件路径",
		Run: func(cmd *cobra.Command, args []string) {
			//通过cmd.flag 获取解析后的flag的值
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			localFilePath, _ := cmd.Flags().GetString("本地文件路径")
			c := webdav_init()
			upload(c, localFilePath, webdavFilePath)

		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "列出文件",
		Long:  "列出文件及目录",
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			c := webdav_init()
			listfile(c, webdavFilePath)
		},
	}

	var delCmd = &cobra.Command{
		Use:   "del",
		Short: "删除文件",
		Long:  "删除文件",
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			c := webdav_init()
			delete(c, webdavFilePath)
		},
	}

	var mkdirCmd = &cobra.Command{
		Use:   "mkdir",
		Short: "建立文件夹",
		Long:  "建立文件夹 支持多层建立",
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			c := webdav_init()
			mkdir(c, webdavFilePath)
		},
	}

	var downCmd = &cobra.Command{
		Use:   "download",
		Short: "下载文件",
		Long:  "下载文件",
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			localFilePath, _ := cmd.Flags().GetString("本地文件路径")
			c := webdav_init()
			download(c, webdavFilePath, localFilePath)
		},
	}
	var init_url = &cobra.Command{
		Use:   "init",
		Short: "初始化url配置",
		Long:  "初始化url配置",
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("访问的url")
			init_url(url)
		},
	}

	var localFilePath string
	var webdavFilePath string

	//添加子命令 upload
	rootCmd.AddCommand(uploadCmd)
	//添加命令下的flag
	uploadCmd.Flags().StringVarP(&localFilePath, "本地文件路径", "f", "/", "本地文件路径")
	uploadCmd.MarkFlagRequired("本地文件路径")

	uploadCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "/", "webdave文件路径")
	uploadCmd.MarkFlagRequired("webdav文件路径")
	uploadCmd.Flags().SortFlags = false

	//添加子命令 list
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "/", "webdave文件路径")
	listCmd.MarkFlagRequired("webdav文件路径")
	listCmd.Flags().SortFlags = false

	//添加子命令 del
	rootCmd.AddCommand(delCmd)
	delCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "", "webdave文件路径")
	delCmd.MarkFlagRequired("webdav文件路径")
	delCmd.Flags().SortFlags = false

	rootCmd.AddCommand(mkdirCmd)
	mkdirCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "", "webdave文件路径")
	mkdirCmd.MarkFlagRequired("webdav文件路径")
	mkdirCmd.Flags().SortFlags = false

	//添加子命令 upload
	rootCmd.AddCommand(downCmd)
	//添加命令下的flag
	downCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "/", "webdave文件路径")
	downCmd.MarkFlagRequired("webdav文件路径")
	downCmd.Flags().StringVarP(&localFilePath, "本地文件路径", "f", "/", "本地文件路径")
	downCmd.MarkFlagRequired("本地文件路径")
	downCmd.Flags().SortFlags = false

	rootCmd.AddCommand(init_url)
	var url string
	init_url.Flags().StringVarP(&url, "访问的url", "u", "https://192.168.31.175", "访问的url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
