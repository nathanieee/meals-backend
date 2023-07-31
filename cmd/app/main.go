package main

import (
	"project-skbackend/configs"
	"project-skbackend/internal/apps"
)

func main() {
	cfg := configs.GetInstance()
	apps.Run(cfg)
}
