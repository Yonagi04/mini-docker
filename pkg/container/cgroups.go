//go:build linux
// +build linux

package container

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	cgroupRoot = "/sys/fs/cgroup"
)

type CgroupManager struct {
	Path string
}

func NewCgroupManager(containerID string) *CgroupManager {
	return &CgroupManager{
		Path: filepath.Join(cgroupRoot, "mini-docker", containerID),
	}
}

func (cm *CgroupManager) Set(container *Container) error {
	// 创建cgroup目录
	if err := os.MkdirAll(cm.Path, 0755); err != nil {
		return fmt.Errorf("创建cgroup目录失败: %v", err)
	}

	// 设置CPU限制
	if container.CPU > 0 {
		if err := cm.setCPU(container.CPU); err != nil {
			return err
		}
	}

	// 设置内存限制
	if container.Memory > 0 {
		if err := cm.setMemory(container.Memory); err != nil {
			return err
		}
	}

	// 将进程PID写入cgroup.procs
	if err := cm.setPID(container.PID); err != nil {
		return err
	}

	return nil
}

func (cm *CgroupManager) setCPU(cpu float64) error {
	// 设置CPU配额
	cpuQuota := int64(cpu * 100000) // 转换为微秒
	if err := os.WriteFile(
		filepath.Join(cm.Path, "cpu.cfs_quota_us"),
		[]byte(strconv.FormatInt(cpuQuota, 10)),
		0644,
	); err != nil {
		return fmt.Errorf("设置CPU配额失败: %v", err)
	}

	// 设置CPU周期
	if err := os.WriteFile(
		filepath.Join(cm.Path, "cpu.cfs_period_us"),
		[]byte("100000"),
		0644,
	); err != nil {
		return fmt.Errorf("设置CPU周期失败: %v", err)
	}

	return nil
}

func (cm *CgroupManager) setMemory(memory int64) error {
	// 设置内存限制
	if err := os.WriteFile(
		filepath.Join(cm.Path, "memory.limit_in_bytes"),
		[]byte(strconv.FormatInt(memory, 10)),
		0644,
	); err != nil {
		return fmt.Errorf("设置内存限制失败: %v", err)
	}

	return nil
}

func (cm *CgroupManager) setPID(pid int) error {
	// 将进程PID写入cgroup.procs
	if err := os.WriteFile(
		filepath.Join(cm.Path, "cgroup.procs"),
		[]byte(strconv.Itoa(pid)),
		0644,
	); err != nil {
		return fmt.Errorf("设置进程PID失败: %v", err)
	}

	return nil
}

func (cm *CgroupManager) Remove() error {
	// 删除cgroup目录
	if err := os.RemoveAll(cm.Path); err != nil {
		return fmt.Errorf("删除cgroup目录失败: %v", err)
	}
	return nil
}
