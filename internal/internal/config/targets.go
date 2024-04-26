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
