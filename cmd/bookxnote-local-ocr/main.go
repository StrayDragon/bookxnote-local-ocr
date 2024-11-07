package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
)

const usageTemplate = `使用方法: %[1]s <命令> [可选参数]

命令:
  server    启动本地转发服务器 (需要管理员权限: 监听443端口)
  cert      证书管理工具
    -force  强制重新生成证书
    -clean  清理并卸载证书

示例:
  %[1]s cert          检查并安装证书(首次安装前需要执行一次, 安装后仅需直接运行 server 命令)
  %[1]s cert -force   强制重新生成证书
  %[1]s cert -clean   清理并卸载证书
  %[1]s server        启动本地转发服务器

`

type App struct {
	executablePath string
	executableDir  string
	binaryName     string
}

func NewApp() (*App, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("无法获取可执行文件路径 > %w", err)
	}

	return &App{
		executablePath: execPath,
		executableDir:  filepath.Dir(execPath),
		binaryName:     filepath.Base(execPath),
	}, nil
}

func (app *App) printUsage() {
	fmt.Printf(usageTemplate, app.binaryName)
}

func runBinary(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (app *App) runBinary(name string, args ...string) error {
	subCmdPath := filepath.Join(app.executableDir, name)
	var newName string
	var newArgs []string
	if runtime.GOOS != "windows" {
		newName = "sudo"
		newArgs = append([]string{subCmdPath}, args...)
	} else {
		newName = subCmdPath
		newArgs = args
	}
	if err := runBinary(newName, newArgs...); err != nil {
		return fmt.Errorf("执行失败 > %w", err)
	}
	return nil
}

func (app *App) ensureCertificates() error {
	certDir := settings.GetCertDir()
	certFiles := []string{"cert.pem", "key.pem"}

	isExist := true
	for _, file := range certFiles {
		if _, err := os.Stat(filepath.Join(certDir, file)); os.IsNotExist(err) {
			isExist = false
			break
		}
	}
	if !isExist {
		if err := app.runBinary("certgen", os.Args[2:]...); err != nil {
			return fmt.Errorf("证书生成失败 > %w", err)
		}
	}
	return nil
}

func (app *App) runServer() error {
	if err := app.ensureCertificates(); err != nil {
		return fmt.Errorf("证书检查失败 > %w", err)
	}

	log.Println("启动本地服务器...")
	if err := app.runBinary("server", os.Args[2:]...); err != nil {
		return fmt.Errorf("服务器启动失败 > %w", err)
	}
	return nil
}

func (app *App) runCertManager() error {
	return app.runBinary("certgen", os.Args[2:]...)
}

func (app *App) Run() error {
	if len(os.Args) < 2 {
		app.printUsage()
		return nil
	}

	switch os.Args[1] {
	case "server":
		return app.runServer()
	case "cert":
		return app.runCertManager()
	default:
		app.printUsage()
		return fmt.Errorf("无效的命令")
	}
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("初始化失败 > %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("运行失败 > %v", err)
	}
}
