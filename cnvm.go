package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 常用常量
var (
	PWD = ""
)

// 定义命令参数
const (
	HELP        = "help"
	USE         = "use"
	LS          = "ls"
	VERSION     = "version"
	UNINSTALL   = "uninstall"
	NODEVERSION = "node-version"
)

func init() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	PWD = strings.Replace(path, "cnvm.exe", "", -1)
}

func main() {

	if len(os.Args) < 2 {
		help()
		return
	}

	switch os.Args[1] {
	case HELP:
		help()
	case USE:
		use()
	case LS:
	case VERSION:
		fmt.Println("cnvm version v1.0.0")
	case "-v":
		fmt.Println("cnvm version v1.0.0")
	case NODEVERSION:
		selectNodeVersion()
	case "node-v":
		selectNodeVersion()
	case UNINSTALL:
		uninstall()
	default:
		help()
	}
}

// 输出命令列表.
func help() {
	fmt.Printf(`
Usage:
  cnvm [flags]
  cnvm [command]

Available Commands:
  use                       Use any the local already exists of Node.js version
  ls                        Show all [local] [remote] Node.js version
  uninstall                 Uninstall local Node.js version and npm
  update                    Update Node.js latest version
  node-version              Show [global] [latest] Node.js version
  version                   Print cnvm version number
	`)
}

// 卸载nodejs版本
func uninstall() {
	err := os.RemoveAll(PWD + os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 切换版本
func use() {
	if !PathExists(PWD + os.Args[2] + "/node.exe") {
		resp, err := http.Get("https://npm.taobao.org/mirrors/node/v" + os.Args[2] + "/win-x64/node.exe")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		content := resp.Body
		if resp.StatusCode == 404 {
			resp, err = http.Get("https://npm.taobao.org/mirrors/node/v" + os.Args[2] + "/node.exe")
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			content = resp.Body
			if resp.StatusCode == 404 {
				fmt.Println("can not found node version.")
				return
			}
		}

		// 查看存放该版本nodejs的文件夹是否存在.
		if !PathExists(PWD + os.Args[2]) {
			os.Mkdir(PWD+os.Args[2], 777)
		}

		// 将下载来了node二进制保存.
		f, err := os.Create(PWD + os.Args[2] + "/node.exe")
		if err != nil {
			panic(err)
		}
		io.Copy(f, content)
	}

	// 复制对应的版本文件到对应目录.
	_, err := CopyFile(PWD+"node.exe", PWD+os.Args[2]+"/node.exe")
	if err != nil {
		fmt.Println(err)
	}
}

// 判断文件是否存在或者该版本nodejs二进制是否存在.
func PathExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return true
}

// 复制文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// 查询nodejs版本
func selectNodeVersion() {
	cmd := exec.Command("node", "-v")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
