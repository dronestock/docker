package config

type Project struct {
	// 仓库地址
	Remote string `default:"${REMOTE=${DRONE_REMOTE_URL=https://github.com/dronestock/docker}}" json:"remote,omitempty"`
	// 镜像链接
	// nolint:lll
	Link string `default:"${LINK=${PLUGIN_REPO_LINK=${DRONE_REPO_LINK=https://github.com/dronestock/docker}}}" json:"link,omitempty"`
}
