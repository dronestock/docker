package config

type Platforms []*Platform

func (p *Platforms) Qemu() (qemu bool) {
	for _, platform := range *p {
		qemu = platform.Qemu()
		if qemu {
			break
		}
	}

	return
}
