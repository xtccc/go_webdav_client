package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/language"
)

var Des_map = map[string]map[language.Tag]string{
	"init_cmd": {
		language.English:           "Please use the init command to initialize the configuration url",
		language.SimplifiedChinese: "请使用init 命令初始化配置url",
	},
	"root_cmd_long": {
		language.English:           "go_webdav_client is an upload and download tool",
		language.SimplifiedChinese: "go_webdav_client 是一个上传下载工具",
	},
	"upload_cmd_long": {
		language.English:           "Uploading files requires two parameters: local file path and webdav file path",
		language.SimplifiedChinese: "上传文件 需要两个参数：本地文件路径 和 webdav文件路径",
	},

	"list_cmd_long": {
		language.English:           "List files and directories",
		language.SimplifiedChinese: "列出文件及目录",
	},

	"del_cmd_long": {
		language.English:           "Delete files/delete folders support recursion",
		language.SimplifiedChinese: "删除文件/删除文件夹 支持递归",
	},
	"mkdir_cmd_long": {
		language.English:           "Create folders to support multi-layer creation",
		language.SimplifiedChinese: "建立文件夹 支持多层建立",
	},
	"download_cmd_long": {
		language.English:           "Downloading a file requires two parameters: local file path and webdav file path",
		language.SimplifiedChinese: "下载文件 需要两个参数：本地文件路径 和 webdav文件路径",
	},
	"init_cmd_long": {
		language.English:           "Initializing url configuration requires remote https link (webdav)",
		language.SimplifiedChinese: "初始化url配置 需要 远程https 链接(webdav)",
	},

	"init_url_flag_des": {
		language.English:           "The accessed url needs to be a webdav service",
		language.SimplifiedChinese: "访问的url,需要是webdav服务",
	},

	"local_file_path_flag_des": {
		language.English:           "local file path",
		language.SimplifiedChinese: "本地文件路径",
	},

	"remote_file_path_flag_des": {
		language.English:           "webdav file path",
		language.SimplifiedChinese: "webdav文件路径",
	},
}

func init_help_str(Des map[string]map[language.Tag]string, help_cmd string) string {
	// 检测用户语言
	langEnv := os.Getenv("LANG")
	langTag := cleanLangTag(langEnv)
	lang, _, err := language.ParseAcceptLanguage(langTag)
	if err != nil {
		fmt.Println(err)
	}
	// 获取对应语言的长描述
	longDesc, exists := Des[help_cmd][lang[0]]
	if !exists {
		longDesc = Des[help_cmd][language.English] // 默认使用英语
	}
	return longDesc
}

// cleanLangTag 函数用于清理语言标签
func cleanLangTag(tag string) string {
	tag = strings.Split(tag, ".")[0]        // 去掉 .UTF-8 等后缀
	tag = strings.ReplaceAll(tag, "_", "-") // 替换 _ 为 -
	if tag == "zh-CN" {
		tag = "zh-Hans"
	}
	return tag
}
