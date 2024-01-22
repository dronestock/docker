package main

import (
	"github.com/dronestock/docker/internal/core"
	"github.com/dronestock/drone"
)

func main() {
	drone.New(core.New).Boot()
}
