package step

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/check"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/rand"
	"github.com/goexl/log"
)

type Boost struct {
	base   *drone.Base
	docker *config.Docker
	config *config.Boost
	logger log.Logger
}

func NewBoost(base *drone.Base, docker *config.Docker, config *config.Boost, logger log.Logger) *Boost {
	return &Boost{
		base:   base,
		docker: docker,
		config: config,
		logger: logger,
	}
}

func (b *Boost) Runnable() bool {
	return nil != b.config && b.config.Enabled
}

func (b *Boost) Run(ctx *context.Context) (err error) {
	dockerfile := b.docker.Dockerfile
	b.docker.Dockerfile = fmt.Sprintf("%s.Dockerfile", rand.New().String().Build().Generate())

	if file, oe := os.Open(dockerfile); nil != oe {
		err = oe
	} else if pe := b.process(ctx, file); nil != pe {
		err = pe
	}

	return
}

func (b *Boost) process(_ *context.Context, file *os.File) (err error) {
	defer func() {
		b.base.Cleanup().Name("清理加速Dockerfile文件").File(b.docker.Dockerfile).Build()
		_ = file.Close()
	}()

	content := new(strings.Builder)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if line, he := b.line(scanner.Text()); nil != he {
			err = he
		} else {
			content.WriteString(line)
			content.WriteString("\n")
		}

		if nil != err {
			return
		}
	}

	fields := gox.Fields[any]{
		field.New("content", content.String()),
		field.New("filename", b.docker.Dockerfile),
	}
	// 写入新文件
	if we := os.WriteFile(b.docker.Dockerfile, []byte(content.String()), os.ModePerm); nil != we {
		err = we
		b.logger.Warn("写入新的Dockerfile出错", fields.Add(field.Error(we))...)
	} else {
		b.logger.Debug("写入新的Dockerfile成功", fields...)
	}

	return
}

func (b *Boost) line(line string) (new string, err error) {
	if !strings.HasPrefix(line, constant.Comment) && strings.Contains(line, constant.From) {
		new, err = b.from(line)
	} else {
		new = line
	}

	return
}

func (b *Boost) from(line string) (new string, err error) {
	if uri, params, pe := b.parse(line); nil != pe {
		err = pe
	} else if b.check(uri.Host) { // 是可以被加速的镜像
		new = fmt.Sprintf("%s %s%s%s", constant.From, b.mirror(uri.Host), uri.Path, gox.If("" != params, params))
	} else if !strings.Contains(uri.Host, constant.Common) { // 是中央仓库的镜像
		host := gox.Ift("" == uri.Path, constant.Library+constant.Slash+uri.Host, uri.Host)
		new = fmt.Sprintf("%s %s/%s%s%s", constant.From, b.config.Mirror, host, uri.Path, gox.If("" != params, params))
	} else {
		new = line
	}

	return
}

func (b *Boost) parse(content string) (uri *url.URL, params string, err error) {
	content = strings.ReplaceAll(content, constant.From, "")
	if strings.Contains(content, constant.Colon) {
		index := strings.Index(content, constant.Colon)
		original := content
		content = original[:index]
		params = original[index:]
	}
	// 组成一个正确的地址去解析
	content = fmt.Sprintf("https://%s", strings.TrimSpace(content))
	uri, err = url.Parse(content)

	return
}

func (b *Boost) check(key string) bool {
	return check.New().Any().String(key).Items(b.config.Hosts...).Contains().Check()
}

func (b *Boost) mirror(host string) (final string) {
	second := host[:strings.Index(host, constant.Common)]
	if constant.DockerProxy == b.config.Mirror {
		final = fmt.Sprintf("%s.%s", second, constant.DockerProxy)
	}

	return
}
