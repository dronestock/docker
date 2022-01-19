package main

import (
	`github.com/dronestock/drone`
)

type plugin struct {
	config *config

	exe               string
	daemonExe         string
	outsideDockerfile string

	daemonSuccessMark string
	loginSuccessMark  string
}

func newPlugin() drone.Plugin {
	return &plugin{
		config: &config{
			defaultMirrors: []string{
				`https://ustc-edu-cn.mirror.aliyuncs.com`,
				`https://mirror.baidubce.com`,
				`https://hub.daocloud.io`,
				`https://mirror.ccs.tencentyun.com`,
			},
		},
		exe:               `/usr/bin/docker`,
		daemonExe:         `/usr/bin/dockerd`,
		outsideDockerfile: `/var/run/docker.sock`,

		daemonSuccessMark: `API listen on /var/run/docker.sock`,
		loginSuccessMark:  `Login Succeeded`,
	}
}

func (p *plugin) Configuration() drone.Configuration {
	return p.config
}

func (p *plugin) Steps() []*drone.Step {
	return []*drone.Step{
		drone.NewStep(p.daemon, drone.Name(`启动守护进程`)),
		drone.NewStep(p.info, drone.Name(`查看Docker信息`)),
		drone.NewStep(p.login, drone.Name(`登录仓库`)),
		drone.NewStep(p.build, drone.Name(`编译镜像`)),
		drone.NewStep(p.push, drone.Name(`推送镜像`)),
	}
}
