package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

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
	ExecutablePath string
	ExecutableDir  string
	BinaryName     string
	HostsManager   *libhosty.HostsFile
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
		ExecutablePath: execPath,
		ExecutableDir:  filepath.Dir(execPath),
		BinaryName:     filepath.Base(execPath),
		HostsManager:   hostsManager,
	}, nil
}

func (app *App) PrintUsage() {
	fmt.Printf(usageTemplate, app.BinaryName)
}

func RunBinary(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (app *App) RunBinary(name string, args ...string) (*exec.Cmd, error) {
	subCmdPath := filepath.Join(app.ExecutableDir, name)
	return RunBinary(subCmdPath, args...)
}

func (app *App) EnsureCertificates() error {
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
		var args []string
		if len(os.Args) > 2 {
			args = os.Args[2:]
		}

		cmd, err := app.RunBinary("certgen", args...)
		if err != nil {
			return fmt.Errorf("证书生成失败 > %w", err)
		}
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("证书生成失败 > %w", err)
		}
	}
	return nil
}

func (app *App) FindHostsLine() (*libhosty.HostsFileLine, error) {
	lines := app.HostsManager.GetHostsFileLines()
	for _, line := range lines {
		if line.Comment == targetHostsComment &&
			line.Address.String() == targetHostsIP &&
			strings.ToLower(line.Hostnames[0]) == targetHostsFQDN {
			return line, nil
		}
	}
	return nil, errors.New("未找到hosts记录")
}

func (app *App) EnsureHosts() error {
	var err error
	if _, err = app.FindHostsLine(); err == nil {
		log.Println("已存在hosts记录, 无须再次配置")
		return nil
	}
	_, _, err = app.HostsManager.AddHostsFileLineRaw(
		targetHostsIP,
		targetHostsFQDN,
		targetHostsComment,
	)
	if err != nil {
		return fmt.Errorf("添加hosts记录失败 > %w", err)
	}

	if err := app.HostsManager.SaveHostsFile(); err != nil {
		return fmt.Errorf("保存hosts文件失败 > %w", err)
	}

	log.Println("hosts配置完成")
	return nil
}

func (app *App) CleanHosts() error {
	line, err := app.FindHostsLine()
	if err != nil {
		return fmt.Errorf("未找到hosts记录 > %w", err)
	}

	app.HostsManager.RemoveHostsFileLineByRow(line.Number)

	if err := app.HostsManager.SaveHostsFile(); err != nil {
		return fmt.Errorf("保存hosts文件失败 > %w", err)
	}

	log.Println("hosts清理完成")
	return nil
}

func (app *App) RunServer() error {
	if err := app.EnsureHosts(); err != nil {
		return fmt.Errorf("hosts配置失败 > %w", err)
	}

	if err := app.EnsureCertificates(); err != nil {
		return fmt.Errorf("证书检查失败 > %w", err)
	}

	log.Println("启动本地服务器...")

	// Only use additional arguments if they exist
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	cmd, err := app.RunBinary("server", args...)
	if err != nil {
		return fmt.Errorf("服务器启动失败 > %w", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	// Handle signals
	go func() {
		select {
		case <-c:
			log.Println("接收到终止信号，正在关闭服务器...")
			if cmd.Process != nil {
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					log.Printf("无法发送终止信号到服务器: %v", err)
					if err := cmd.Process.Kill(); err != nil {
						log.Printf("无法强制终止服务器: %v", err)
					}
				}
			}
		case <-done:
			// Server already exited, do nothing
		}
	}()

	log.Println("服务器已启动，按Ctrl+C终止...")

	// Wait for server to exit
	err = <-done
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			log.Println("服务器已退出")
		} else {
			return fmt.Errorf("服务器运行出错 > %w", err)
		}
	}

	// Reset signal handling
	signal.Reset(os.Interrupt, syscall.SIGTERM)

	return nil
}

func (app *App) ConfirmAction(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(response)
	return response == "y"
}

func (app *App) RunInstall() error {
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "-cert":
			var args []string
			if len(os.Args) > 3 && os.Args[3] == "-force" {
				args = []string{"-force"}
			}
			cmd, err := app.RunBinary("certgen", args...)
			if err != nil {
				return fmt.Errorf("证书安装失败 > %w", err)
			}
			return cmd.Wait()
		case "-hosts":
			return app.EnsureHosts()
		default:
			return fmt.Errorf("无效的安装参数")
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "install" {
		if !app.ConfirmAction("即将安装所有配置, 请确保使用管理员运行或设置过权限! 是否继续?") {
			return fmt.Errorf("用户取消安装")
		}
	}

	if err := app.EnsureCertificates(); err != nil {
		return fmt.Errorf("证书安装失败 > %w", err)
	}

	if err := app.EnsureHosts(); err != nil {
		return fmt.Errorf("hosts配置失败 > %w", err)
	}

	log.Println("安装完成")
	return nil
}

func (app *App) RunUninstall() error {
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "-cert":
			certDir := settings.GetCertDir()
			for _, f := range []string{"cert.pem", "key.pem"} {
				path := filepath.Join(certDir, f)
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("无法删除证书文件 %s > %w", f, err)
				}
			}
			log.Println("证书文件已删除")
			return nil
		case "-hosts":
			return app.CleanHosts()
		default:
			return fmt.Errorf("无效的卸载参数")
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "uninstall" {
		if !app.ConfirmAction("即将卸载所有配置, 是否继续?") {
			return fmt.Errorf("用户取消卸载")
		}
	}

	if err := app.CleanHosts(); err != nil {
		log.Printf("警告: 无法清理hosts文件: %v", err)
	}

	certDir := settings.GetCertDir()
	for _, f := range []string{"cert.pem", "key.pem"} {
		path := filepath.Join(certDir, f)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			log.Printf("警告: 无法删除证书文件 %s: %v", f, err)
		}
	}

	log.Println("卸载完成")
	return nil
}

func (app *App) Run() error {
	if len(os.Args) < 2 {
		app.PrintUsage()
		return nil
	}

	cmd := os.Args[1]
	switch cmd {
	case "server":
		return app.RunServer()
	case "install":
		return app.RunInstall()
	case "uninstall":
		return app.RunUninstall()
	default:
		app.PrintUsage()
		return fmt.Errorf("未知命令: %s", cmd)
	}
}
