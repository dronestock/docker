package main

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/goexl/gox"
	"github.com/goexl/gox/check"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/rand"
)

type stepBoost struct {
	*plugin
}

func newBoostStep(plugin *plugin) *stepBoost {
	return &stepBoost{
		plugin: plugin,
	}
}

func (b *stepBoost) Runnable() bool {
	return b.boostEnabled()
}

func (b *stepBoost) Run(ctx context.Context) (err error) {
	dockerfile := b.Dockerfile
	b.Dockerfile = fmt.Sprintf("%s.Dockerfile", rand.New().String().Generate())

	if file, oe := os.Open(dockerfile); nil != oe {
		err = oe
	} else if pe := b.process(ctx, file); nil != pe {
		err = pe
	}

	return
}

func (b *stepBoost) process(_ context.Context, file *os.File) (err error) {
	defer func() {
		b.Cleanup().Name("清理加速Dockerfile文件").File(b.Dockerfile).Build()
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
		field.New("filename", b.Dockerfile),
	}
	// 写入新文件
	if we := os.WriteFile(b.Dockerfile, []byte(content.String()), os.ModePerm); nil != we {
		err = we
		b.Warn("写入新的Dockerfile出错", fields.Connect(field.Error(we))...)
	} else {
		b.Debug("写入新的Dockerfile成功", fields...)
	}

	return
}

func (b *stepBoost) line(line string) (new string, err error) {
	if strings.Contains(line, from) {
		new, err = b.from(line)
	} else {
		new = line
	}

	return
}

func (b *stepBoost) from(line string) (new string, err error) {
	if uri, params, pe := b.parse(line); nil != pe {
		err = pe
	} else if b.check(uri.Host) { // 是可以被加速的镜像
		new = fmt.Sprintf("%s %s%s%s", from, b.mirror(uri.Host), uri.Path, gox.If("" != params, params))
	} else if !strings.Contains(uri.Host, common) { // 是中央仓库的镜像
		new = fmt.Sprintf("%s %s/%s%s%s", from, b.Boost.Mirror, uri.Host, uri.Path, gox.If("" != params, params))
	}

	return
}

func (b *stepBoost) parse(content string) (uri *url.URL, params string, err error) {
	content = strings.ReplaceAll(content, from, "")
	if strings.Contains(content, colon) {
		index := strings.Index(content, colon)
		original := content
		content = original[:index]
		params = original[index:]
	}
	// 组成一个正确的地址去解析
	content = fmt.Sprintf("https://%s", strings.TrimSpace(content))
	uri, err = url.Parse(content)

	return
}

func (b *stepBoost) check(key string) bool {
	return check.New().Any().String(key).Items(b.Boost.Hosts...).Contains().Check()
}

func (b *stepBoost) mirror(host string) (final string) {
	second := host[:strings.Index(host, common)]
	if dockerProxy == b.Boost.Mirror {
		final = fmt.Sprintf("%s.%s", second, dockerProxy)
	}

	return
}
