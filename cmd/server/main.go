package main

import (
	"fmt"
	"web_server/configs"
)

func main() {
	cfg, _ := configs.LoadConfig(".")
	fmt.Println(cfg.DBDriver)
}
