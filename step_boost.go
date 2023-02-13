package main

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/goexl/gox"
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
	b.Dockerfile = rand.New().String().Generate()

	if file, oe := os.Open(dockerfile); nil != oe {
		err = oe
	} else if pe := b.process(ctx, file); nil != pe {
		err = pe
	}

	return
}

func (b *stepBoost) process(_ context.Context, file *os.File) (err error) {
	defer func() {
		_ = file.Close()
	}()

	content := new(strings.Builder)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if line, he := b.has(scanner.Text()); nil != he {
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
		b.Warn("写入新的Dockerfile成功", fields...)
	}

	return
}

func (b *stepBoost) has(content string) (new string, err error) {
	lowercase := strings.ToLower(content)
	if !strings.Contains(lowercase, from) {
		new = content
	} else if uri, pe := url.Parse(strings.ReplaceAll(lowercase, from, "")); nil != pe {
		err = pe
	} else if "" == uri.Host {
		new = fmt.Sprintf("FROM %s/%s", b.Boost.Mirror, strings.TrimSpace(uri.Path))
	} else if b.check(uri.Host) {
		new = fmt.Sprintf("%s/%s/%s", b.Boost.Mirror, uri.Host, strings.TrimSpace(uri.Path))
	} else {
		new = content
	}

	return
}

func (b *stepBoost) check(key string) (has bool) {
	has = false
	for _, check := range b.Boost.Hosts {
		if strings.Contains(key, check) {
			has = true
		}

		if has {
			break
		}
	}

	return
}
