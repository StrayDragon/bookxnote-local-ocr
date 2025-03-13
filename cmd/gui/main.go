//go:generate fyne bundle -o resources.go ../../artifact/icon.png
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/straydragon/bookxnote-local-ocr/internal/gui"
)

func main() {
	app, err := gui.NewGUIApp(resourceIconPng)
	if err != nil {
		log.Fatalf("初始化GUI应用程序失败: %v", err)
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
