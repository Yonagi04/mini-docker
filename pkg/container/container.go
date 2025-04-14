package container

import (
	"fmt"
	"sync"
	"time"
)

var (
	idGenerator IDGenerator
	initOnce    sync.Once
)

func initIDGenerator() error {
	var err error
	initOnce.Do(func() {
		idGenerator, err = NewSnowflakeIDGenerator(1, 1)
	})
	return err
}

type Container struct {
	ID        string
	PID       int
	Command   string
	Status    string
	Created   time.Time
	Name      string
	Image     string   // 使用的镜像
	WorkDir   string   // 工作目录
	Env       []string // 环境变量
	Volumes   []string // 挂载的卷
	Network   string   // 网络模式
	IPAddress string   // IP地址
	Ports     []string // 端口映射
	Memory    int64    // 内存限制（字节）
	CPU       float64  // CPU限制（核数）
}

type ContainerOpts struct {
	Commands []string
	Tty      bool
	Name     string
	Image    string
	WorkDir  string
	Env      []string
	Volumes  []string
	Network  string
	Ports    []string
	Memory   int64
	CPU      float64
}

func NewContainer(opts *ContainerOpts) (*Container, error) {
	// 确保ID生成器已经初始化
	if err := initIDGenerator(); err != nil {
		return nil, fmt.Errorf("初始化ID生成器失败: %v", err)
	}

	id := idGenerator.NextID()

	name := opts.Name
	// 如果名字为空，使用ID的前8位作为名称
	if name == "" {
		name = id[:8]
	}

	return &Container{
		ID:        id,
		Command:   opts.Commands[0],
		Status:    "created",
		Created:   time.Now(),
		Name:      name,
		Image:     opts.Image,
		WorkDir:   opts.WorkDir,
		Env:       opts.Env,
		Volumes:   opts.Volumes,
		Network:   opts.Network,
		IPAddress: "",
		Ports:     opts.Ports,
		Memory:    opts.Memory,
		CPU:       opts.CPU,
	}, nil
}

func (c *Container) Start() error {
	if c.Status != "created" && c.Status != "stopped" {
		return fmt.Errorf("容器状态错误: %s", c.Status)
	}
	c.Status = "running"
	return nil
}

func (c *Container) Pause() error {
	if c.Status != "running" {
		return fmt.Errorf("容器状态错误: %s", c.Status)
	}
	c.Status = "paused"
	return nil
}

func (c *Container) Resume() error {
	if c.Status != "paused" {
		return fmt.Errorf("容器状态错误: %s", c.Status)
	}
	c.Status = "running"
	return nil
}

func (c *Container) Remove() error {
	if c.Status == "running" {
		return fmt.Errorf("无法删除运行中的容器")
	}
	c.Status = "removed"
	return nil
}

func (c *Container) String() string {
	return fmt.Sprintf("Container(ID: %s, Name: %s, Status: %s, Created: %s, Image: %s)",
		c.ID, c.Name, c.Status, c.Created.Format(time.RFC3339), c.Image)
}
