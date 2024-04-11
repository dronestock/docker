package config

type Docker struct {
	// 主机
	Host string `default:"${HOST=/var/run/docker.sock}" json:"host,omitempty" validate:"required"`
	// 端口
	Port uint `default:"${PORT}" json:"port,omitempty" validate:"max=65535"`
	// 协议
	Protocol string `default:"${PROTOCOL=${PROTO=unix}}" json:"protocol,omitempty" validate:"oneof=ssh tcp unix"`
	// 用户名
	// 因为USERNAME环境变量在Docker镜像里面已经被使用
	// 所以先用USER环境变量去接收参数
	// 如果USER没设置而USERNAME被设置，那么接收到的值是正确的，因为环境变量是可以被覆盖的
	Username string `default:"${USER=${USERNAME}}" json:"username,omitempty" validate:"required_if=Protocol ssh"`
	// 密码
	Password string `default:"${PASSWORD}" json:"password,omitempty" validate:"required_if=Protocol ssh Key ''"`
	// 密钥
	Key string `default:"${KEY}" json:"key,omitempty" validate:"required_if=Protocol ssh Password ''"`
	// 启动成功后的标志
	Mark string `default:"${MARK=API listen on /var/run/docker.sock}" json:"mark,omitempty"`

	// 镜像列表
	Mirrors []string `default:"${MIRRORS}" json:"mirrors,omitempty"`
	// 启用实验性功能
	Experimental bool `default:"${EXPERIMENTAL=true}" json:"experimental,omitempty"`
	// 精减镜像层数
	Squash bool `default:"${SQUASH=true}" json:"squash,omitempty"`
	// 压缩镜像
	Compress bool `default:"${COMPRESS=true}" json:"compress,omitempty"`
	// 标签列表
	Labels []string `default:"${LABELS}" json:"labels,omitempty"`

	// 数据目录
	Data string `default:"${DATA=/var/lib/docker}" json:"data,omitempty"`
	// 驱动
	Driver string `default:"${DRIVER}" json:"driver,omitempty"`

	// 仓库
	Repository string `default:"${REPOSITORY}" json:"repository,omitempty"`
	// 执行程序
	Exe string `default:"${EXE=/usr/bin/docker}" json:"exe,omitempty"`
	// 执行程序
	Daemon string `default:"${DAEMON=/usr/bin/dockerd}" json:"exe,omitempty"`
}
