package main

import (
	"context"
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

func (i *stepInfo) Run(_ context.Context) error {
	return i.Command(exe).Args("info").Dir(i.context()).Exec()
}
