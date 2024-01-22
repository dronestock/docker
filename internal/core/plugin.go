package core

import (
	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/docker/internal/step"
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type plugin struct {
	drone.Base
	// 核心配置
	config.Docker `default:"${DOCKER}"`
	// 加速
	config.Boost `default:"${BOOST}"`
	// 项目
	config.Project `default:"${PROJECT}"`

	// 仓库
	Registry *config.Registry `default:"${REGISTRY}"`
	// 仓库列表
	Registries []*config.Registry `default:"${REGISTRIES}"`
}

func New() drone.Plugin {
	return new(plugin)
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Steps() drone.Steps {
	return drone.Steps{
		drone.NewStep(step.NewSsh(&p.Base, &p.Docker, p.Logger)).Name("授权").Build(),
		drone.NewStep(step.NewBoost(&p.Base, &p.Docker, &p.Boost, p.Logger)).Name("加速").Build(),
		drone.NewStep(step.NewDaemon(&p.Base, &p.Docker, p.Logger)).Name("守护").Build(),
		drone.NewStep(step.NewLogin(&p.Base, &p.Docker, p.Registries, p.Logger)).Name("登录").Build(),
		drone.NewStep(step.NewBuild(&p.Base, &p.Docker, &p.Project, p.Registries, p.Logger)).Name("编译").Build(),
		drone.NewStep(step.NewPush(&p.Base, &p.Docker, p.Registries, p.Logger)).Name("推送").Build(),
	}
}

func (p *plugin) Setup() (err error) {
	if nil != p.Registry {
		p.Registries = append(p.Registries, p.Registry)
	}

	return
}

func (p *plugin) Fields() gox.Fields[any] {
	return gox.Fields[any]{
		field.New("dockerfile", p.Dockerfile),
		field.New("context", p.Context),
		field.New("host", p.Host),
		field.New("mirrors", p.Mirrors),
		field.New("tag", p.Tag),
		field.New("tag.auto", p.Auto),
		field.New("name", p.Name),

		field.New("experimental", p.Experimental),
		field.New("squash", p.Squash),
		field.New("compress", p.Compress),
		field.New("labels", p.Labels),

		field.New("remote", p.Remote),
		field.New("link", p.Link),

		field.New("registries", p.Registries),
		field.New("repository", p.Repository),
	}
}
