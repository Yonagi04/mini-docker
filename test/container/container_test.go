package container_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/yonagi04/mini-docker/pkg/container"
)

func getTestWorkDir() string {
	if runtime.GOOS == "windows" {
		return "C:\\Windows\\Temp"
	}
	return "/tmp"
}

func getTestEnv() []string {
	if runtime.GOOS == "windows" {
		return []string{
			"Path=C:\\Windows\\System32;C:\\Windows",
			"SYSTEMROOT=C:\\Windows",
			"TEMP=C:\\Windows\\Temp",
		}
	}
	return []string{
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"TERM=xterm",
	}
}

func getTestCommand(command ...string) []string {
	if runtime.GOOS == "windows" {
		return append([]string{"C:\\Windows\\System32\\cmd.exe", "/c"}, command...)
	}
	return command
}

// 通用测试用例，不依赖于特定操作系统
func TestContainerBasics(t *testing.T) {
	tests := []struct {
		name    string
		opts    *container.ContainerOpts
		wantErr bool
	}{
		{
			name: "空命令测试",
			opts: &container.ContainerOpts{
				Commands: []string{},
				Tty:      true,
			},
			wantErr: true,
		},
		{
			name: "基本命令测试",
			opts: &container.ContainerOpts{
				Commands: getTestCommand("echo", "hello"),
				Tty:      true,
				WorkDir:  getTestWorkDir(),
				Env:      getTestEnv(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := container.NewContainer(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if c.ID == "" {
					t.Error("NewContainer() 返回的容器ID为空")
				}
				if c.Status != "created" {
					t.Errorf("NewContainer() 初始状态错误，got = %v, want = created", c.Status)
				}
				if time.Since(c.Created) > time.Second {
					t.Error("NewContainer() 创建时间不正确")
				}
			}
		})
	}
}
