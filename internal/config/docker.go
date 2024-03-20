package config

type Docker struct {
	// 配置文件
	Dockerfile string `default:"${DOCKERFILE=Dockerfile}" json:"dockerfile,omitempty" validate:"required"`
	// 上下文
	Context string `default:"${CONTEXT=.}" json:"context,omitempty"`

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
	// 标签
	Tag string `default:"${TAG=${DRONE_TAG=0.0.${DRONE_BUILD_NUMBER}}}" json:"tag,omitempty" validate:"required"`
	// 前缀
	Prefix string `default:"${PREFIX}" json:"prefix,omitempty"`
	// 后缀
	Suffix string `default:"${SUFFIX}" json:"suffix,omitempty"`
	// 自动标签
	Auto bool `default:"${AUTO=true}" json:"auto,omitempty"`
	// 名称
	Name string `default:"${NAME=${DRONE_COMMIT_SHA=latest}}" json:"name,omitempty"`

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
}
