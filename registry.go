package main

type registry struct {
	// 仓库地址
	Hostname string `default:"docker.io" json:"hostname" validate:"required,hostname"`
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
	// 是否必须成功
	Required bool `json:"required"`
}
