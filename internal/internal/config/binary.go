package config

type Binary struct {
	// 执行程序
	Docker string `default:"${DOCKER=/usr/bin/docker}" json:"docker,omitempty"`
	// 执行程序
	Daemon string `default:"${DAEMON=/usr/bin/dockerd}" json:"daemon,omitempty"`
}
