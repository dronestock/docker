package step

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
)

const sshConfigFormatter = `Host *
  IgnoreUnknown UseKeychain
  UseKeychain yes
  AddKeysToAgent yes
  StrictHostKeyChecking=no
  IdentityFile %s
`

type Ssh struct {
	base   *drone.Base
	docker *config.Docker
	logger log.Logger
}

func NewSsh(base *drone.Base, docker *config.Docker, logger log.Logger) *Ssh {
	return &Ssh{
		base:   base,
		docker: docker,
		logger: logger,
	}
}

func (s *Ssh) Runnable() bool {
	return "" != s.docker.Key || "" != s.docker.Password
}

func (s *Ssh) Run(_ *context.Context) (err error) {
	home := filepath.Join(os.Getenv(constant.HomeEnv), constant.SshHome)
	keyfile := filepath.Join(home, constant.SshKeyFilename)
	configFile := filepath.Join(home, constant.SshConfigDir)
	if me := s.makeSSHHome(home); nil != me { // 创建主目录
		err = me
	} else if we := s.writeSSHKey(keyfile); nil != we { // 写入密钥文件
		err = we
	} else if ce := s.writeSSHConfig(configFile, keyfile); nil != ce { // 写入配置文件
		err = ce
	} else { // 设置环境变量
		_ = os.Setenv(constant.DockerHost, s.host())
	}

	return
}

func (s *Ssh) makeSSHHome(home string) (err error) {
	homeField := field.New("home", home)
	if err = os.MkdirAll(home, os.ModePerm); nil != err {
		s.logger.Error("创建SSH目录出错", homeField, field.Error(err))
	}

	return
}

func (s *Ssh) writeSSHKey(keyfile string) (err error) {
	key := s.docker.Key
	keyfileField := field.New("keyfile", keyfile)
	// 必须以换行符结束
	if !strings.HasSuffix(key, "\n") {
		key = fmt.Sprintf("%s\n", key)
	}

	if err = os.WriteFile(keyfile, []byte(key), constant.DefaultFilePerm); nil != err {
		s.logger.Error("写入密钥文件出错", keyfileField, field.Error(err))
	}

	return
}

func (s *Ssh) writeSSHConfig(configFile string, keyfile string) (err error) {
	configFileField := field.New("file", configFile)
	if err = os.WriteFile(configFile, []byte(fmt.Sprintf(sshConfigFormatter, keyfile)), constant.DefaultFilePerm); nil != err {
		s.logger.Error("写入SSH配置文件出错", configFileField, field.Error(err))
	}

	return
}

func (s *Ssh) host() string {
	builder := new(strings.Builder)
	builder.WriteString(s.docker.Protocol)
	builder.WriteString(constant.Colon)
	builder.WriteString(constant.Slash)
	builder.WriteString(constant.Slash)
	if "" != s.docker.Username {
		builder.WriteString(s.docker.Username)
		builder.WriteString(constant.Dollar)
	}
	builder.WriteString(s.docker.Host)
	if 0 != s.docker.Port {
		builder.WriteString(constant.Colon)
		builder.WriteString(gox.ToString(s.docker.Port))
	}

	return builder.String()
}
