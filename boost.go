package main

type boost struct {
	// 是否开户
	Enabled *bool `default:"true" json:"enabled"`
	// 加速服务器
	Mirror string `default:"dockerproxy.com" json:"mirror"`
	// 可被加速的地址
	Hosts []string `default:"['ghcr.io', 'gcr.io', 'k8s.gcr.io', 'quay.io']" json:"hosts"`
}
