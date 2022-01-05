package main

import (
	`fmt`
	`io/ioutil`
	`os`
	`path/filepath`
	`strings`

	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

const sshConfig = `Host *
  IgnoreUnknown UseKeychain
  UseKeychain yes
  AddKeysToAgent yes
  StrictHostKeyChecking=no
  IdentityFile %s
`

func ssh(conf *config, logger simaqian.Logger) (err error) {
	if `` == conf.SSHKey {
		return
	}

	home := filepath.Join(os.Getenv(`HOME`), `.ssh`)
	keyfile := filepath.Join(home, `id_rsa`)
	configFile := filepath.Join(home, `config`)
	if err = makeSSHHome(home, logger); nil != err {
		return
	}
	if err = writeSSHKey(keyfile, conf.SSHKey, logger); nil != err {
		return
	}
	err = writeSSHConfig(configFile, keyfile, logger)

	return
}

func makeSSHHome(home string, logger simaqian.Logger) (err error) {
	homeField := field.String(`home`, home)
	if err = os.MkdirAll(home, os.ModePerm); nil != err {
		logger.Error(`创建SSH目录出错`, homeField, field.Error(err))
	} else {
		logger.Info(`创建SSH目录成功`, homeField)
	}

	return
}

func writeSSHKey(keyfile string, key string, logger simaqian.Logger) (err error) {
	keyfileField := field.String(`keyfile`, keyfile)
	// 必须以换行符结束
	if !strings.HasSuffix(key, `\n`) {
		key = fmt.Sprintf(`%s\n`, key)
	}

	if err = ioutil.WriteFile(keyfile, []byte(key), 0600); nil != err {
		logger.Error(`写入密钥文件出错`, keyfileField, field.Error(err))
	} else {
		logger.Info(`写入密钥文件成功`, keyfileField)
	}

	return
}

func writeSSHConfig(configFile string, keyfile string, logger simaqian.Logger) (err error) {
	configFileField := field.String(`config.file`, configFile)
	if err = ioutil.WriteFile(configFile, []byte(fmt.Sprintf(sshConfig, keyfile)), 0600); nil != err {
		logger.Error(`写入SSH配置文件出错`, configFileField, field.Error(err))
	} else {
		logger.Info(`写入SSH配置文件成功`, configFileField)
	}

	return
}
