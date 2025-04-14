//go:build linux
// +build linux

package container

import (
	"os/exec"
	"syscall"
)

func (c *Container) setupNamespaces(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // 隔离主机名
			syscall.CLONE_NEWPID | // 隔离进程
			syscall.CLONE_NEWNS | // 隔离挂载点
			syscall.CLONE_NEWNET | // 隔离网络
			syscall.CLONE_NEWIPC | // 隔离IPC
			syscall.CLONE_NEWUSER, // 隔离用户
	}
	return nil
}
