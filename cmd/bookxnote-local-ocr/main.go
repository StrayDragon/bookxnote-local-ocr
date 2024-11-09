package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/areYouLazy/libhosty"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
)

const (
	targetHostsComment = "bookxnote-local-ocr"
	targetHostsIP      = "127.0.0.1"
	targetHostsFQDN    = "aip.baidubce.com"
)

const usageTemplate = `BookXNote 本地OCR方案
项目地址: https://github.com/straydragon/bookxnote-local-ocr

使用方法: %[1]s <命令> [可选参数]

命令:
  server              启动本地转发服务器 (需要管理员权限: 监听443端口)
  install             安装所有必需配置 (需要管理员权限)
    -cert [-force]    安装证书 [-force: 强制重新生成]
    -hosts            配置hosts
  uninstall           卸载所有配置 (需要管理员权限)
    -cert             仅卸载证书
    -hosts            仅清理hosts

示例:
  %[1]s install           安装所有必需配置
  %[1]s install -cert     仅安装证书
  %[1]s install -hosts    仅配置hosts
  %[1]s uninstall        卸载所有配置
  %[1]s server           启动本地转发服务器

`

type App struct {
	executablePath string
	executableDir  string
	binaryName     string
	hostsManager   *libhosty.HostsFile
}

func NewApp() (*App, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("无法获取可执行文件路径 > %w", err)
	}

	hostsManager, err := libhosty.Init()
	if err != nil {
		return nil, fmt.Errorf("初始化hosts管理器失败 > %w", err)
	}

	return &App{
		executablePath: execPath,
		executableDir:  filepath.Dir(execPath),
		binaryName:     filepath.Base(execPath),
		hostsManager:   hostsManager,
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
	// var newName string
	// var newArgs []string
	// if runtime.GOOS != "windows" {
	// 	newName = "sudo"
	// 	newArgs = append([]string{subCmdPath}, args...)
	// } else {
	// 	newName = subCmdPath
	// 	newArgs = args
	// }
	if err := runBinary(subCmdPath, args...); err != nil {
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

func (app *App) findHostsLine() (*libhosty.HostsFileLine, error) {
	lines := app.hostsManager.GetHostsFileLines()
	for _, line := range lines {
		if line.Comment == targetHostsComment &&
			line.Address.String() == targetHostsIP &&
			strings.ToLower(line.Hostnames[0]) == targetHostsFQDN {
			return line, nil
		}
	}
	return nil, errors.New("未找到hosts记录")
}

func (app *App) ensureHosts() error {
	var err error
	if _, err = app.findHostsLine(); err == nil {
		log.Println("已存在hosts记录, 无须再次配置")
		return nil
	}
	_, _, err = app.hostsManager.AddHostsFileLineRaw(
		targetHostsIP,
		targetHostsFQDN,
		targetHostsComment,
	)
	if err != nil {
		return fmt.Errorf("添加hosts记录失败 > %w", err)
	}

	if err := app.hostsManager.SaveHostsFile(); err != nil {
		return fmt.Errorf("保存hosts文件失败 > %w", err)
	}

	log.Println("hosts配置完成")
	return nil
}

func (app *App) cleanHosts() error {
	line, err := app.findHostsLine()
	if err != nil {
		return fmt.Errorf("未找到hosts记录 > %w", err)
	}

	app.hostsManager.RemoveHostsFileLineByRow(line.Number)

	if err := app.hostsManager.SaveHostsFile(); err != nil {
		return fmt.Errorf("保存hosts文件失败 > %w", err)
	}

	log.Println("hosts清理完成")
	return nil
}

func (app *App) runServer() error {
	if err := app.ensureHosts(); err != nil {
		return fmt.Errorf("hosts配置失败 > %w", err)
	}

	if err := app.ensureCertificates(); err != nil {
		return fmt.Errorf("证书检查失败 > %w", err)
	}

	log.Println("启动本地服务器...")
	if err := app.runBinary("server", os.Args[2:]...); err != nil {
		return fmt.Errorf("服务器启动失败 > %w", err)
	}
	return nil
}

func (app *App) confirmAction(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(response)
	return response == "y"
}

func (app *App) runInstall() error {
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "-cert":
			var args []string
			if len(os.Args) > 3 && os.Args[3] == "-force" {
				args = []string{"-force"}
			}
			return app.runBinary("certgen", args...)
		case "-hosts":
			return app.ensureHosts()
		default:
			return fmt.Errorf("无效的安装参数")
		}
	}

	if !app.confirmAction("即将安装所有配置, 请确保使用管理员权限! 是否继续?") {
		return fmt.Errorf("用户取消安装")
	}

	if err := app.ensureHosts(); err != nil {
		return fmt.Errorf("hosts配置失败 > %w", err)
	}

	if err := app.runBinary("certgen"); err != nil {
		return fmt.Errorf("证书安装失败 > %w", err)
	}

	log.Println("所有配置安装完成")
	return nil
}

func (app *App) runUninstall() error {
	// 解析参数
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "-cert":
			// 转发到原来的cert -clean命令
			return app.runBinary("certgen", "-clean")
		case "-hosts":
			return app.cleanHosts()
		default:
			return fmt.Errorf("无效的卸载参数")
		}
	}

	// 卸载全部配置
	if !app.confirmAction("即将卸载所有配置, 请确保使用管理员权限~ 是否继续?") {
		return fmt.Errorf("用户取消卸载")
	}

	if err := app.runBinary("certgen", "-clean"); err != nil {
		log.Printf("警告: 证书卸载失败 > %v", err)
	}

	if err := app.cleanHosts(); err != nil {
		log.Printf("警告: hosts清理失败 > %v", err)
	}

	log.Println("所有配置已卸载")
	return nil
}

func checkAdminPrivileges() error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("net", "session")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("需要管理员权限运行! 请右键点击程序，选择「以管理员身份运行」")
		}
	} else if runtime.GOOS == "linux" {
		if os.Geteuid() != 0 {
			return fmt.Errorf("需要root权限运行! 请使用 sudo 运行此程序")
		}
	}
	return nil
}

func (app *App) Run() error {
	if err := checkAdminPrivileges(); err != nil {
		return err
	}

	if len(os.Args) < 2 {
		app.printUsage()
		return nil
	}

	switch os.Args[1] {
	case "server":
		return app.runServer()
	case "install":
		return app.runInstall()
	case "uninstall":
		return app.runUninstall()
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
