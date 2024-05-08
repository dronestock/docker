package config

import (
	"strings"

	"github.com/dronestock/docker/internal/internal/constant"
)

type Targets []*Target

func (t *Targets) Runnable(registries *Registries, docker *Docker) (runnable bool) {
	for _, target := range *t {
		if target.Pushable(registries, docker) {
			runnable = true
		}
		if runnable {
			break
		}
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
		all = append(all, target.PlatformArgument())
	}
	if 0 != len(all) {
		platforms = strings.Join(all, constant.Comma)
	}

	return
}

func (t *Targets) Binfmt(docker *Docker) (fmt bool) {
	if "/var/run/docker.sock" == docker.Host {
		fmt = t.qemu()
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
