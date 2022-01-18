package main

import (
	`github.com/dronestock/drone`
)

func main() {
	panic(drone.Bootstrap(newPlugin()))
}
