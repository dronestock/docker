package main

import (
	`fmt`
	`os`
	`strings`

	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/mengpo`
	`github.com/storezhang/validatorx`
)

type config struct {
	// 目录
	Folder string `default:"${PLUGIN_FOLDER=${FOLDER=.}}" validate:"required"`
	// 镜像列表
	Mirrors []string `default:"${PLUGIN_MIRRORS=${MIRRORS}}"`
	// 标签
	Tag string `default:"${PLUGIN_TAG=${TAG}}"`

	// 是否启用默认值
	Defaults bool `default:"${PLUGIN_DEFAULTS=${DEFAULTS=true}}"`
	// 是否显示调试信息
	Verbose bool `default:"${PLUGIN_VERBOSE=${VERBOSE=false}}"`

	envs []string
}

func (c *config) Fields() gox.Fields {
	return []gox.Field{
		field.String(`remote`, c.Remote),
		field.String(`folder`, c.Folder),
		field.String(`branch`, c.Branch),
		field.String(`tag`, c.Tag),
		field.String(`author`, c.Author),
		field.String(`email`, c.Email),
		field.String(`message`, c.Message),

		field.Bool(`clear`, c.Clear),
		field.Bool(`verbose`, c.Verbose),
	}
}

func (c *config) load() (err error) {
	if err = mengpo.Set(c); nil != err {
		return
	}
	if err = validatorx.Struct(c); nil != err {
		return
	}
	c.envs = make([]string, 0)

	return
}

func (c *config) addEnvs(envs ...*env) {
	for _, _env := range envs {
		c.envs = append(c.envs, fmt.Sprintf(`%s=%s`, _env.key, _env.value))
	}
}

func (c *config) pull() bool {
	return `1` == os.Getenv(`DRONE_STEP_NUMBER`) || `pull` == c.Mode
}

func (c *config) fastGithub() bool {
	return strings.HasPrefix(c.Remote, `https://github.com`) || strings.HasPrefix(c.Remote, `http://github.com`)
}

func (c *config) gitForce() (force string) {
	if c.Force {
		force = `--force`
	}

	return
}

func (c *config) checkout() (checkout string) {
	if `` != c.Commit {
		checkout = c.Commit
	} else {
		checkout = c.Branch
	}

	return
}
