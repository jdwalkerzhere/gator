package main

import (
	"fmt"
	"github.com/jdwalkerzhere/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	cfg.SetUser("Jesse")
	cfg, err = config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.CurrentUser)
}
