package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	UPDATE      = "update"
)

func init() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	PWD = strings.Replace(path, "cnvm.exe", "", -1)
}

func main() {
	// 如果命令少于2位则匹配不到对应的逻辑.
	if len(os.Args) < 2 {
		help()
		return
	}

	switch os.Args[1] {
	case HELP:
		help()
	case USE:
		use()
	case UPDATE:
		Update()
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

// 切换版本
func use() {
	if !PathExists(PWD + os.Args[2] + "/node.exe") {
		num, err := strconv.ParseInt(strings.Split(os.Args[2], ".")[0], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		str := ""
		if num > int64(4) {
			str = "/win-x64"
		}

		resp, err := http.Get("https://npm.taobao.org/mirrors/node/v" + os.Args[2] + str + "/node.exe")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			fmt.Println("can not found node version.")
			return
		}

		// 查看存放该版本nodejs的文件夹是否存在.
		if !PathExists(PWD + os.Args[2]) {
			os.Mkdir(PWD+os.Args[2], 777)
		}

		// 将下载来了node二进制保存.
		f, err := os.Create(PWD + os.Args[2] + "/node.exe")
		if err != nil {
			panic(err)
			return
		}
		io.Copy(f, resp.Body)
	}

	// 复制对应的版本文件到对应目录.
	_, err := CopyFile(PWD+"node.exe", PWD+os.Args[2]+"/node.exe")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("use node version " + os.Args[2] + " successfuly")
}

// 更新到最近版本的Node
func Update() {
	resp, err := http.Get("https://npm.taobao.org/mirrors/node/latest/win-x64/node.exe")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		fmt.Println("can not found node version.")
		return
	}
	// 将下载来了node二进制保存.
	f, err := os.Create(PWD + "/node.exe")
	if err != nil {
		panic(err)
		return
	}
	io.Copy(f, resp.Body)

	fmt.Println("use node version latest successfuly")
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

// 卸载nodejs版本
func uninstall() {
	err := os.RemoveAll(PWD + os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("remove node version " + os.Args[2] + " successfuly")
}
