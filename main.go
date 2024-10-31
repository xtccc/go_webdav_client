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
		fmt.Println(init_help_str(Des_map, "init_cmd"))
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
		Use: "go_webdav_client",
		//表示当前命令的完整描述。
		Long: init_help_str(Des_map, "root_cmd_long"),

		//Run属性是一个函数，当执行命令时会调用此函数。
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "upload file",
		Long:  init_help_str(Des_map, "upload_cmd_long"),
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
		Short: "list file",
		Long:  init_help_str(Des_map, "list_cmd_long"),
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			c := webdav_init()
			listfile(c, webdavFilePath)
		},
	}

	var delCmd = &cobra.Command{
		Use:   "del",
		Short: "delete file",
		Long:  init_help_str(Des_map, "del_cmd_long"),
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			c := webdav_init()
			delete(c, webdavFilePath)
		},
	}

	var mkdirCmd = &cobra.Command{
		Use:   "mkdir",
		Short: "mkdir",
		Long:  init_help_str(Des_map, "mkdir_cmd_long"),
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			c := webdav_init()
			mkdir(c, webdavFilePath)
		},
	}

	var downCmd = &cobra.Command{
		Use:   "download",
		Short: "download file",
		Long:  init_help_str(Des_map, "download_cmd_long"),
		Run: func(cmd *cobra.Command, args []string) {
			webdavFilePath, _ := cmd.Flags().GetString("webdav文件路径")
			localFilePath, _ := cmd.Flags().GetString("本地文件路径")
			c := webdav_init()
			download(c, webdavFilePath, localFilePath)
		},
	}
	var init_url = &cobra.Command{
		Use:   "init",
		Short: "init",
		Long:  init_help_str(Des_map, "init_cmd_long"),
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
	uploadCmd.Flags().StringVarP(&localFilePath, "本地文件路径", "f", "/", init_help_str(Des_map, "local_file_path_flag_des"))
	uploadCmd.MarkFlagRequired("本地文件路径")

	uploadCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "/", init_help_str(Des_map, "remote_file_path_flag_des"))
	uploadCmd.MarkFlagRequired("webdav文件路径")
	uploadCmd.Flags().SortFlags = false

	//添加子命令 list
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "/", init_help_str(Des_map, "remote_file_path_flag_des"))
	listCmd.MarkFlagRequired("webdav文件路径")
	listCmd.Flags().SortFlags = false

	//添加子命令 del
	rootCmd.AddCommand(delCmd)
	delCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "", init_help_str(Des_map, "remote_file_path_flag_des"))
	delCmd.MarkFlagRequired("webdav文件路径")
	delCmd.Flags().SortFlags = false

	rootCmd.AddCommand(mkdirCmd)
	mkdirCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "", init_help_str(Des_map, "remote_file_path_flag_des"))
	mkdirCmd.MarkFlagRequired("webdav文件路径")
	mkdirCmd.Flags().SortFlags = false

	//添加子命令 upload
	rootCmd.AddCommand(downCmd)
	//添加命令下的flag
	downCmd.Flags().StringVarP(&webdavFilePath, "webdav文件路径", "w", "/", init_help_str(Des_map, "remote_file_path_flag_des"))
	downCmd.MarkFlagRequired("webdav文件路径")
	downCmd.Flags().StringVarP(&localFilePath, "本地文件路径", "f", "/", init_help_str(Des_map, "local_file_path_flag_des"))
	downCmd.MarkFlagRequired("本地文件路径")
	downCmd.Flags().SortFlags = false

	rootCmd.AddCommand(init_url)
	var url string
	init_url.Flags().StringVarP(&url, "访问的url", "u", "https://192.168.31.175", init_help_str(Des_map, "init_url_flag_des"))
	init_url.MarkFlagRequired("访问的url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
