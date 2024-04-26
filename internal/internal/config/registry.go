package config

import (
	"strings"
)

type Registry struct {
	// 仓库地址
	Hostname string `default:"docker.io" json:"hostname" validate:"required,hostname"`
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
	// 是否必须成功
	Required bool `json:"required"`
	// 登录成功标志
	Mark string `default:"Login Succeeded" json:"mark"`

	Name string `default:"未设置"`
}

func (r *Registry) Nickname() (nickname string) {
	switch strings.TrimSpace(r.Hostname) {
	case "ccr.ccs.tencentyun.com":
		nickname = "腾讯云"
	case "docker.io":
		nickname = "中央库"
	default:
		nickname = r.Name
	}

	return
}
