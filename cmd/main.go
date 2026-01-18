package main

import (
	"fmt"

	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/pkg/parse"
)

func main() {
	path := parse.FlagInit()

	cfg := config.MustLoadConfig(path)
	fmt.Println(cfg)
}
