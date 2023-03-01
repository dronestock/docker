package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goexl/gox/field"
)

const sshConfigFormatter = `Host *
  IgnoreUnknown UseKeychain
  UseKeychain yes
  AddKeysToAgent yes
  StrictHostKeyChecking=no
  IdentityFile %s
`

type stepSsh struct {
	*plugin
}

func newSshStep(plugin *plugin) *stepSsh {
	return &stepSsh{
		plugin: plugin,
	}
}

func (s *stepSsh) Runnable() bool {
	return "" != s.Key
}

func (s *stepSsh) Run(_ context.Context) (err error) {
	home := filepath.Join(os.Getenv(homeEnv), sshHome)
	keyfile := filepath.Join(home, sshKeyFilename)
	configFile := filepath.Join(home, sshConfigDir)
	if err = s.makeSSHHome(home); nil != err {
		return
	}
	if err = s.writeSSHKey(keyfile); nil != err {
		return
	}
	err = s.writeSSHConfig(configFile, keyfile)

	return
}

func (s *stepSsh) makeSSHHome(home string) (err error) {
	homeField := field.New("home", home)
	if err = os.MkdirAll(home, os.ModePerm); nil != err {
		s.Error("创建SSH目录出错", homeField, field.Error(err))
	}

	return
}

func (s *stepSsh) writeSSHKey(keyfile string) (err error) {
	key := s.Key
	keyfileField := field.New("keyfile", keyfile)
	// 必须以换行符结束
	if !strings.HasSuffix(key, "\n") {
		key = fmt.Sprintf("%s\n", key)
	}

	if err = os.WriteFile(keyfile, []byte(key), defaultFilePerm); nil != err {
		s.Error("写入密钥文件出错", keyfileField, field.Error(err))
	}

	return
}

func (s *stepSsh) writeSSHConfig(configFile string, keyfile string) (err error) {
	configFileField := field.New("file", configFile)
	if err = os.WriteFile(configFile, []byte(fmt.Sprintf(sshConfigFormatter, keyfile)), defaultFilePerm); nil != err {
		s.Error("写入SSH配置文件出错", configFileField, field.Error(err))
	}

	return
}
