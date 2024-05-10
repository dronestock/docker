package internal

import (
	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/step"
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type plugin struct {
	drone.Base

	// 执行程序
	Binary config.Binary `default:"${BINARY}" json:"binary,omitempty"`
	// 核心配置
	Docker config.Docker `default:"${DOCKER}" json:"docker,omitempty"`
	// 加速
	Boost config.Boost `default:"${BOOST}" json:"boost,omitempty"`
	// 项目
	Project config.Project `default:"${PROJECT}" json:"project,omitempty"`

	// 目标
	Target config.Target `default:"${TARGET}" json:"target,omitempty"`
	// 目标列表
	Targets config.Targets `default:"${TARGETS}" json:"targets,omitempty"`

	// 仓库
	Registry *config.Registry `default:"${REGISTRY}" json:"registry,omitempty"`
	// 仓库列表
	Registries config.Registries `default:"${REGISTRIES}" json:"registries,omitempty"`

	docker  *command.Docker
	targets *config.Targets
}

func New() drone.Plugin {
	return new(plugin)
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Steps() drone.Steps {
	return drone.Steps{
		drone.NewStep(step.NewSSH(&p.Base, &p.Docker, p.Logger)).Name("授权").Build(),
		drone.NewStep(step.NewBoost(&p.Base, p.Targets, &p.Boost, p.Logger)).Name("加速").Build(),
		drone.NewStep(step.NewDaemon(p.docker, &p.Docker)).Name("守护").Build(),
		drone.NewStep(step.NewLogin(p.docker, &p.Docker, &p.Registries, p.targets)).Name("登录").Build(),
		drone.NewStep(step.NewSetup(p.docker, &p.Docker, p.targets)).Name("配置").Build(),
		drone.NewStep(step.NewBuild(p.docker, &p.Docker, p.targets, &p.Project, &p.Registries)).Name("编译").Build(),
		drone.NewStep(step.NewPush(p.docker, &p.Docker, p.targets, &p.Registries)).Name("推送").Build(),
	}
}

func (p *plugin) Setup() (err error) {
	if nil != p.Registry { // 合并仓库
		p.Registries = append(p.Registries, p.Registry)
	}
	// 创建执行封装命令
	p.docker = command.NewDocker(&p.Base, &p.Binary)

	return
}

func (p *plugin) Fields() gox.Fields[any] {
	return gox.Fields[any]{
		field.New("targets", p.Targets),
		field.New("docker", p.Docker),
		field.New("registries", p.Registries),
		field.New("project", p.Project),
		field.New("binary", p.Binary),
	}
}

func (p *plugin) Before() (err error) {
	p.Targets = append(p.Targets, &p.Target)
	cache := make(map[string]*config.Target)
	for _, target := range p.Targets {
		tag := target.Tag
		if cached, ok := cache[tag]; ok {
			cached.Platforms = append(cached.Platforms, &target.Platform)
			cached.Platforms = append(cached.Platforms, target.Platforms...)
		} else {
			cache[tag] = target
		}
	}

	targets := make(config.Targets, 0, len(cache))
	for _, target := range cache {
		targets = append(targets, target)
	}
	p.targets = &targets

	return
}
