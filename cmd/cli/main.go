package main

import (
	"fmt"
	"log"
	"os"

	"github.com/straydragon/bookxnote-local-ocr/internal/cli"
)

func main() {
	app, err := cli.NewApp()
	if err != nil {
		log.Fatalf("初始化应用程序失败: %v", err)
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
