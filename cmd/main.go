package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/yonagi04/mini-docker/pkg/container"
)

func getDefaultPath() string {
	if runtime.GOOS == "windows" {
		return "Path=C:\\Windows\\System32;C:\\Windows;C:\\Windows\\System32\\Wbem;C:\\Windows\\System32\\WindowsPowerShell\\v1.0"
	}
	return "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
}

func getDefaultWorkDir() string {
	if runtime.GOOS == "windows" {
		return "C:\\"
	}
	return "/"
}

func getCommand() []string {
	if runtime.GOOS == "windows" {
		return []string{"C:\\Windows\\System32\\cmd.exe"}
	}
	return os.Args[1:]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: mini-docker run [command]")
		if runtime.GOOS == "windows" {
			fmt.Println("例如: mini-docker run cmd")
		} else {
			fmt.Println("例如: mini-docker run /bin/sh")
		}
		os.Exit(1)
	}

	command := getCommand()
	opts := &container.ContainerOpts{
		Commands: command,
		Tty:      true,
		WorkDir:  getDefaultWorkDir(),
		Env:      []string{getDefaultPath()},
	}

	// 创建新容器
	c, err := container.NewContainer(opts)
	if err != nil {
		fmt.Printf("创建容器失败: %v\n", err)
		os.Exit(1)
	}

	// 运行容器
	if err := c.Run(opts); err != nil {
		fmt.Printf("运行容器失败: %v\n", err)
		os.Exit(1)
	}
}
