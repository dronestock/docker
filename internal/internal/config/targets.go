package config

import (
	"strings"

	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/goexl/gox"
)

type Targets []*Target

func (t *Targets) Runnable(registries *Registries, docker *Docker) (runnable bool) {
	buildPushedCount := 0
	for _, target := range *t {
		if target.BuildWithPush(registries, docker) {
			buildPushedCount = buildPushedCount + 1
		}
	}
	if len(*t) == buildPushedCount { // 如果全部都在编译时就推送了，不执行推送步骤
		runnable = false
	} else { // 视仓库数量为准
		runnable = 0 != len(*registries)
	}

	return
}

func (t *Targets) Registries() (registries Registries) {
	registries = make(Registries, 0, len(*t))
	for _, target := range *t {
		registries = append(registries, target.AllRegistries()...)
	}

	return
}

func (t *Targets) Platforms() (platforms string) {
	all := make([]string, 0, len(*t))
	for _, target := range *t {
		all = gox.Ift("" != target.PlatformArgument(), append(all, target.PlatformArgument()), all)
	}
	if 0 != len(all) {
		platforms = strings.Join(all, constant.Comma)
	}

	return
}

func (t *Targets) BinfmtNeedSetup(docker *Docker) (fmt bool) {
	if "/var/run/docker.sock" == docker.Host {
		fmt = t.qemu()
	}

	return
}

func (t *Targets) DriverNeedSetup() (need bool) {
	for _, target := range *t {
		need = strings.Contains(target.PlatformArgument(), constant.Comma)
		if need {
			break
		}
	}

	return
}

func (t *Targets) qemu() (qemu bool) {
	for _, target := range *t {
		qemu = target.Qemu()
		if qemu {
			break
		}
	}

	return
}
