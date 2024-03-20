package main

import (
	"github.com/dronestock/docker/internal"
	"github.com/dronestock/drone"
)

func main() {
	drone.New(internal.New).Boot()
}
