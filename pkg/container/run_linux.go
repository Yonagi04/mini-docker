//go:build linux
// +build linux

package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func (c *Container) Run(opts *ContainerOpts) error {
	// 检查容器状态
	if c.Status != "created" && c.Status != "stopped" {
		return fmt.Errorf("容器状态错误: %s", c.Status)
	}

	// 创建工作目录
	if err := os.MkdirAll(c.WorkDir, 0755); err != nil {
		return fmt.Errorf("创建工作目录失败: %v", err)
	}

	// 设置环境变量
	cmd := exec.Command(opts.Commands[0], opts.Commands[1:]...)
	cmd.Dir = c.WorkDir
	cmd.Env = append(os.Environ(), c.Env...)

	// 设置标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 设置命名空间
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // 隔离主机名
			syscall.CLONE_NEWPID | // 隔离进程
			syscall.CLONE_NEWNS | // 隔离挂载点
			syscall.CLONE_NEWNET | // 隔离网络
			syscall.CLONE_NEWIPC | // 隔离IPC
			syscall.CLONE_NEWUSER, // 隔离用户
	}

	// 启动容器
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动容器失败: %v", err)
	}

	c.PID = cmd.Process.Pid
	c.Status = "running"
	log.Printf("容器 %s 已启动，PID: %d", c.ID, c.PID)

	// 设置资源限制
	cgroupManager := NewCgroupManager(c.ID)
	if err := cgroupManager.Set(c); err != nil {
		log.Printf("警告：设置资源限制失败: %v", err)
	}

	// 设置网络
	if err := c.setupNetwork(); err != nil {
		return fmt.Errorf("设置网络失败: %v", err)
	}

	// 挂载卷
	if err := c.mountVolumes(); err != nil {
		return fmt.Errorf("挂载卷失败: %v", err)
	}

	if opts.Tty {
		return cmd.Wait()
	}
	return nil
}

func (c *Container) Stop() error {
	if c.PID <= 0 {
		return fmt.Errorf("容器未运行")
	}
	if c.Status != "running" {
		return fmt.Errorf("容器状态错误: %s", c.Status)
	}

	// 发送终止信号
	process, err := os.FindProcess(c.PID)
	if err != nil {
		return fmt.Errorf("查找进程失败: %v", err)
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("停止容器失败: %v", err)
	}

	// 等待进程结束
	_, err = process.Wait()
	if err != nil {
		return fmt.Errorf("等待进程结束失败: %v", err)
	}

	// 清理资源限制
	cgroupManager := NewCgroupManager(c.ID)
	if err := cgroupManager.Remove(); err != nil {
		log.Printf("警告：清理资源限制失败: %v", err)
	}

	c.Status = "stopped"
	log.Printf("容器 %s 已停止", c.ID)
	return nil
}

func (c *Container) setupNetwork() error {
	// TODO: 实现网络设置
	return nil
}

func (c *Container) mountVolumes() error {
	for _, volume := range c.Volumes {
		parts := filepath.SplitList(volume)
		if len(parts) != 2 {
			return fmt.Errorf("无效的卷格式: %s", volume)
		}
		source := parts[0]
		target := parts[1]

		// 确保源目录存在
		if err := os.MkdirAll(source, 0755); err != nil {
			return fmt.Errorf("创建源目录失败: %v", err)
		}

		// 确保目标目录存在
		if err := os.MkdirAll(target, 0755); err != nil {
			return fmt.Errorf("创建目标目录失败: %v", err)
		}

		// TODO: 实现挂载
	}
	return nil
}
