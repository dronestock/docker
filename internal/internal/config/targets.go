package config

type Targets []*Target

func (t *Targets) Runnable() (runnable bool) {
	for _, target := range *t {
		if nil != target.Registry || 0 != len(target.Registries) {
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
