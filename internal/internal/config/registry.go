package config

import (
	"strings"

	"github.com/dronestock/docker/internal/internal/constant"
)

type Registry struct {
	// 仓库地址
	Hostname string `default:"docker.io" json:"hostname,omitempty" validate:"required,hostname"`
	// 用户名
	Username string `json:"username,omitempty"`
	// 密码
	Password string `json:"password,omitempty"`
	// 登录成功标志
	Mark string `default:"Login Succeeded" json:"mark,omitempty"`

	Name string `default:"未设置"`
}

func (r *Registry) Nickname() (nickname string) {
	switch strings.TrimSpace(r.Hostname) {
	case "ccr.ccs.tencentyun.com":
		nickname = "腾讯云"
	case constant.DockerIO:
		nickname = "中央库"
	default:
		nickname = r.Name
	}

	return
}

func (r *Registry) DockerHub() bool {
	return constant.DockerIO == r.Hostname
}
