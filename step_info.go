package main

import (
	"context"
	"errors"
	"time"

	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
)

type stepInfo struct {
	*plugin
}

func newInfoStep(plugin *plugin) *stepInfo {
	return &stepInfo{
		plugin: plugin,
	}
}

func (i *stepInfo) Runnable() bool {
	return true
}

func (i *stepInfo) Run(ctx context.Context) (err error) {
	times := 0
	for {
		times++
		if err = i.check(ctx); nil == err || errors.Is(err, context.DeadlineExceeded) {
			break
		} else {
			i.Info("等待Docker启动完成", field.New("times", times), field.New("address", i.address()))
		}
		time.Sleep(5 * time.Second)
	}

	return
}

func (i *stepInfo) check(ctx context.Context) (err error) {
	ia := args.New().Build().Subcommand("info").Build()
	_, err = i.Command(exe).Args(ia).Dir(i.context()).Context(ctx).Build().Exec()

	return
}
