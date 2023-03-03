package main

import (
	"context"
	"fmt"
	"os"
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
	fmt.Println(os.Getenv(dockerHost))
	return i.Command(exe).Args("info").Dir(i.context()).Build().Exec()
}
