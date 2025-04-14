package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

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

func getDefaultEnv() []string {
	if runtime.GOOS == "windows" {
		return []string{
			getDefaultPath(),
			"SYSTEMROOT=C:\\Windows",
			"TEMP=C:\\Windows\\Temp",
			"TMP=C:\\Windows\\Temp",
			"COMSPEC=C:\\Windows\\System32\\cmd.exe",
			"PROMPT=$P$G",
		}
	}
	return []string{
		getDefaultPath(),
		"TERM=xterm",
		"HOME=/root",
		"SHELL=/bin/sh",
	}
}

func getCommand() []string {
	if runtime.GOOS == "windows" {
		// 如果没有指定具体命令，则直接运行cmd
		if len(os.Args) <= 2 {
			return []string{"C:\\Windows\\System32\\cmd.exe"}
		}
		// 如果指定了命令，则使用/c参数执行
		cmd := []string{"C:\\Windows\\System32\\cmd.exe", "/c"}
		cmd = append(cmd, os.Args[2:]...)
		return cmd
	}
	// Linux系统直接使用传入的命令
	return os.Args[1:]
}

func printUsage() {
	fmt.Println("使用方法: mini-docker run [command]")
	if runtime.GOOS == "windows" {
		fmt.Println("例如:")
		fmt.Println("  mini-docker run cmd                    # 运行cmd")
		fmt.Println("  mini-docker run cmd dir                # 执行dir命令")
		fmt.Println("  mini-docker run cmd echo Hello World   # 执行echo命令")
	} else {
		fmt.Println("例如:")
		fmt.Println("  mini-docker run /bin/sh               # 运行shell")
		fmt.Println("  mini-docker run /bin/ls               # 执行ls命令")
		fmt.Println("  mini-docker run /bin/echo Hello World # 执行echo命令")
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := getCommand()
	opts := &container.ContainerOpts{
		Commands: command,
		Tty:      true,
		WorkDir:  getDefaultWorkDir(),
		Env:      getDefaultEnv(),
	}

	// 创建新容器
	c, err := container.NewContainer(opts)
	if err != nil {
		fmt.Printf("创建容器失败: %v\n", err)
		os.Exit(1)
	}

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在单独的goroutine中运行容器
	go func() {
		if err := c.Run(opts); err != nil {
			fmt.Printf("运行容器失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待信号
	<-sigChan
	fmt.Println("\n接收到终止信号，正在关闭容器...")

	// 清理容器资源
	if err := c.Cleanup(); err != nil {
		fmt.Printf("清理容器资源失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("容器已成功关闭")
}
