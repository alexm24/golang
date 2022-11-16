package main

import (
	"github.com/alexm24/golang/internal/app"
)

func main() {
	const configPath = "configs/config.yaml"
	app.App(configPath)
}
