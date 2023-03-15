package main

import (
	"context"

	"github.com/goexl/gox/args"
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

func (i *stepInfo) Run(_ context.Context) (err error) {
	ia := args.New().Build().Subcommand("info").Build()
	_, err = i.Command(exe).Args(ia).Dir(i.context()).Build().Exec()

	return
}
