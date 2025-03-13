package gui

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/straydragon/bookxnote-local-ocr/internal/cli"
	"github.com/straydragon/bookxnote-local-ocr/internal/client/openapi"
)

var serverStopChan = make(chan struct{}, 1)
var updateSystemTray func()

var i18n = map[string]map[string]string{
	"zh": {
		"appId":    "com.straydragon.bookxnote-local-ocr",
		"appTitle": "BookXNote Local OCR",

		"status.stopped":     "状态: 已停止",
		"status.starting":    "状态: 正在启动...",
		"status.running":     "状态: 运行中",
		"status.startFailed": "状态: 启动失败",
		"status.stopping":    "状态: 正在停止...",

		"ui.startServer":     "启动服务器",
		"ui.stopServer":      "停止服务器",
		"ui.installConfig":   "安装配置",
		"ui.uninstallConfig": "卸载配置",
		"ui.show":            "显示",
		"ui.exit":            "退出",

		"ui.postOCR":               "OCR后处理",
		"ui.postOCR.autoFixLayout": "自动整理行",
		"ui.postOCR.translate":     "自动翻译",
		"ui.postOCR.generateNotes": "自动生成笔记",

		"notify.title.server":       "服务器状态",
		"notify.title.error":        "错误",
		"notify.title.installing":   "安装中",
		"notify.title.success":      "成功",
		"notify.title.uninstalling": "卸载中",
		"notify.title.config":       "配置更新",

		"notify.msg.serverStarted":     "服务器已启动",
		"notify.msg.serverStopped":     "服务器已停止",
		"notify.msg.installing":        "正在安装配置...",
		"notify.msg.installComplete":   "配置安装完成",
		"notify.msg.uninstalling":      "正在卸载配置...",
		"notify.msg.uninstallComplete": "配置卸载完成",
		"notify.msg.configUpdated":     "配置已更新",

		"error.serverAlreadyRunning": "服务器已在运行",
		"error.serverNotRunning":     "服务器未运行",
		"error.startFailed":          "启动服务器失败: %v",
		"error.stopFailed":           "停止服务器失败: %v",
		"error.installFailed":        "安装配置失败: %v",
		"error.uninstallFailed":      "卸载配置失败: %v",
		"error.configUpdateFailed":   "更新配置失败: %v",

		"log.serverStarted":  "[BXN Local OCR (GUI)] %v | 服务器已启动\n",
		"log.serverStopped":  "[BXN Local OCR (GUI)] %v | 服务器已停止\n",
		"log.serverStopping": "[BXN Local OCR (GUI)] %v | 已停止服务器...\n",
	},
}

var currentLang = "zh"

func t(key string) string {
	if val, ok := i18n[currentLang][key]; ok {
		return val
	}
	return key
}

type GUIApp struct {
	*cli.App
	fyneApp       fyne.App
	window        fyne.Window
	serverProcess *os.Process
	serverSwitch  *widget.Check
}

func NewGUIApp(resourceIconPng fyne.Resource) (*GUIApp, error) {
	cliApp, err := cli.NewApp()
	if err != nil {
		return nil, err
	}

	fyneApp := app.NewWithID(t("appId"))
	fyneApp.SetIcon(resourceIconPng)

	window := fyneApp.NewWindow(t("appTitle"))

	guiApp := &GUIApp{
		App:     cliApp,
		fyneApp: fyneApp,
		window:  window,
	}

	return guiApp, nil
}

func (g *GUIApp) setupUI() {
	statusLabel := widget.NewLabel(t("status.stopped"))

	g.serverSwitch = widget.NewCheck(t("ui.startServer"), func(isChecked bool) {
		if isChecked {
			g.handleServerStart(statusLabel)
		} else {
			g.handleServerStop(statusLabel)
		}
	})
	g.serverSwitch.SetChecked(g.serverProcess != nil)

	installBtn := widget.NewButton(t("ui.installConfig"), func() {
		go g.handleInstallConfig()
	})

	uninstallBtn := widget.NewButton(t("ui.uninstallConfig"), func() {
		go g.handleUninstallConfig()
	})

	content := container.NewVBox(
		statusLabel,
		g.serverSwitch,
		installBtn,
		uninstallBtn,
	)

	g.window.SetContent(content)
}

func (g *GUIApp) runServerWithContext(ctx context.Context) error {
	if err := g.App.EnsureHosts(); err != nil {
		return fmt.Errorf("hosts配置失败: %w", err)
	}

	if err := g.App.EnsureCertificates(); err != nil {
		return fmt.Errorf("证书检查失败: %w", err)
	}

	serverPath := filepath.Join(g.App.ExecutableDir, "server")
	var args []string // No additional arguments for now

	cmd := exec.Command(serverPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("服务器启动失败: %w", err)
	}

	log.Printf("服务器进程已启动 (PID: %d)", cmd.Process.Pid)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		log.Println("正在终止服务器进程...")

		err := cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Printf("无法发送SIGTERM信号: %v，尝试SIGKILL", err)
			_ = cmd.Process.Kill()
		} else {
			time.Sleep(500 * time.Millisecond)

			if err := cmd.Process.Signal(syscall.Signal(0)); err == nil {
				log.Println("服务器进程未响应，强制终止")
				_ = cmd.Process.Kill()
			}
		}

		<-done
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (g *GUIApp) handleServerStart(statusLabel *widget.Label) {
	if g.serverProcess != nil {
		return
	}

	go func() {
		statusLabel.SetText(t("status.starting"))

		errCh := make(chan error, 1)
		doneCh := make(chan struct{}, 1)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			go func() {
				select {
				case <-serverStopChan:
					cancel()
				case <-ctx.Done():
				}
			}()

			err := g.runServerWithContext(ctx)
			if err != nil {
				select {
				case <-ctx.Done():
				default:
					errCh <- fmt.Errorf("服务器错误: %w", err)
				}
			}
			doneCh <- struct{}{}
		}()

		select {
		case err := <-errCh:
			dialog.ShowError(fmt.Errorf(t("error.startFailed"), err), g.window)
			g.serverSwitch.SetChecked(false)
			statusLabel.SetText(t("status.startFailed"))
			cancel()
			return
		case <-time.After(500 * time.Millisecond):
			g.serverProcess = &os.Process{Pid: os.Getpid()}
			statusLabel.SetText(t("status.running"))
			g.notifyUser(t("notify.title.server"), t("notify.msg.serverStarted"))
			fmt.Printf(t("log.serverStarted"), time.Now().Format("2006/01/02 - 15:04:05"))
		}

		go func() {
			<-doneCh
			g.serverProcess = nil
			g.serverSwitch.SetChecked(false)
			statusLabel.SetText(t("status.stopped"))
			g.notifyUser(t("notify.title.server"), t("notify.msg.serverStopped"))
			fmt.Printf(t("log.serverStopped"), time.Now().Format("2006/01/02 - 15:04:05"))
			cancel()
		}()
	}()
}

func (g *GUIApp) handleServerStop(statusLabel *widget.Label) {
	if g.serverProcess == nil {
		return
	}

	statusLabel.SetText(t("status.stopping"))

	select {
	case serverStopChan <- struct{}{}:
	default:
	}

	statusLabel.SetText(t("status.stopped"))
}

func (g *GUIApp) handleInstallConfig() {
	dialog.ShowInformation(t("notify.title.installing"), t("notify.msg.installing"), g.window)

	err := g.App.RunInstall()
	if err != nil {
		dialog.ShowError(fmt.Errorf(t("error.installFailed"), err), g.window)
	} else {
		dialog.ShowInformation(t("notify.title.success"), t("notify.msg.installComplete"), g.window)
	}
}

func (g *GUIApp) handleUninstallConfig() {
	dialog.ShowInformation(t("notify.title.uninstalling"), t("notify.msg.uninstalling"), g.window)

	err := g.App.RunUninstall()
	if err != nil {
		dialog.ShowError(fmt.Errorf(t("error.uninstallFailed"), err), g.window)
	} else {
		dialog.ShowInformation(t("notify.title.success"), t("notify.msg.uninstallComplete"), g.window)
	}
}

func (g *GUIApp) notifyUser(title, content string) {
	notification := fyne.NewNotification(title, content)
	g.fyneApp.SendNotification(notification)
}

func (g *GUIApp) startServer() error {
	if g.serverProcess != nil {
		return errors.New(t("error.serverAlreadyRunning"))
	}

	errCh := make(chan error, 1)
	doneCh := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		go func() {
			select {
			case <-serverStopChan:
				cancel()
			case <-ctx.Done():
			}
		}()

		err := g.runServerWithContext(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
			default:
				errCh <- fmt.Errorf("服务器错误: %w", err)
			}
		}
		doneCh <- struct{}{}
	}()

	select {
	case err := <-errCh:
		cancel()
		return fmt.Errorf(t("error.startFailed"), err)
	case <-time.After(500 * time.Millisecond):
		g.serverProcess = &os.Process{Pid: os.Getpid()}
		g.notifyUser(t("notify.title.server"), t("notify.msg.serverStarted"))
		fmt.Printf(t("log.serverStarted"), time.Now().Format("2006/01/02 - 15:04:05"))

		if updateSystemTray != nil {
			updateSystemTray()
		}
	}

	if g.serverSwitch != nil {
		g.serverSwitch.SetChecked(true)
	}

	go func() {
		<-doneCh
		g.serverProcess = nil
		if g.serverSwitch != nil {
			g.serverSwitch.SetChecked(false)
		}
		fmt.Printf(t("log.serverStopped"), time.Now().Format("2006/01/02 - 15:04:05"))

		if updateSystemTray != nil {
			updateSystemTray()
		}

		cancel()
	}()

	return nil
}

func (g *GUIApp) stopServer() error {
	if g.serverProcess == nil {
		return errors.New(t("error.serverNotRunning"))
	}

	select {
	case serverStopChan <- struct{}{}:
	default:
	}

	g.notifyUser(t("notify.title.server"), t("notify.msg.serverStopped"))
	fmt.Printf(t("log.serverStopping"), time.Now().Format("2006/01/02 - 15:04:05"))

	g.serverProcess = nil

	if g.serverSwitch != nil {
		g.serverSwitch.SetChecked(false)
	}

	if updateSystemTray != nil {
		updateSystemTray()
	}

	return nil
}

func (g *GUIApp) setupTray() {
	desktopApp, ok := g.fyneApp.(desktop.App)
	if !ok {
		return
	}

	var toggleServerItem *fyne.MenuItem

	updateMenu := func() {
		if g.serverProcess == nil {
			toggleServerItem.Label = t("ui.startServer")
		} else {
			toggleServerItem.Label = t("ui.stopServer")
		}

		menu := createTrayMenu(g, toggleServerItem)
		desktopApp.SetSystemTrayMenu(menu)
	}

	updateSystemTray = updateMenu

	toggleServerItem = fyne.NewMenuItem(t("ui.startServer"), func() {
		if g.serverProcess == nil {
			if err := g.startServer(); err != nil {
				g.notifyUser(t("notify.title.error"), err.Error())
				return
			}
		} else {
			if err := g.stopServer(); err != nil {
				g.notifyUser(t("notify.title.error"), err.Error())
				return
			}
		}
		updateMenu()
	})

	menu := createTrayMenu(g, toggleServerItem)
	desktopApp.SetSystemTrayMenu(menu)

	go g.updateMenuPeriodically(desktopApp, toggleServerItem, updateMenu)
}

func createTrayMenu(g *GUIApp, toggleServerItem *fyne.MenuItem) *fyne.Menu {
	showItem := fyne.NewMenuItem(t("ui.show"), func() {
		g.window.Show()
	})

	autoFixLayoutItem := fyne.NewMenuItem(t("ui.postOCR.autoFixLayout"), func() {
		go g.togglePostOCRFeature("after_ocr.auto_fix_content.enabled")
	})

	translateItem := fyne.NewMenuItem(t("ui.postOCR.translate"), func() {
		go g.togglePostOCRFeature("after_ocr.translate.enabled")
	})

	generateNotesItem := fyne.NewMenuItem(t("ui.postOCR.generateNotes"), func() {
		go g.togglePostOCRFeature("after_ocr.generate_by_llm.enabled")
	})

	exitItem := fyne.NewMenuItem(t("ui.exit"), func() {
		g.fyneApp.Quit()
	})

	menuItems := []*fyne.MenuItem{toggleServerItem, showItem}

	if g.serverProcess != nil {
		menuItems = append(
			menuItems,
			fyne.NewMenuItemSeparator(),
			autoFixLayoutItem,
			translateItem,
			generateNotesItem,
		)
	}

	menuItems = append(menuItems, fyne.NewMenuItemSeparator(), exitItem)

	return fyne.NewMenu("", menuItems...)
}

func (g *GUIApp) updateMenuPeriodically(_ desktop.App, _ *fyne.MenuItem, updateMenu func()) {
	updateMenu()

	var lastServerState bool
	for {
		time.Sleep(1 * time.Second)
		if g.fyneApp == nil {
			return
		}

		currentServerState := g.serverProcess != nil
		if currentServerState != lastServerState {
			updateMenu()
			lastServerState = currentServerState
		}
	}
}

func (g *GUIApp) Run() error {
	g.setupUI()
	g.setupTray()

	g.window.Resize(fyne.NewSize(300, 200))
	g.window.CenterOnScreen()
	g.window.SetCloseIntercept(func() {
		g.window.Hide()
	})

	g.window.ShowAndRun()
	return nil
}

func RunGUI(resourceIconPng fyne.Resource) error {
	guiApp, err := NewGUIApp(resourceIconPng)
	if err != nil {
		return err
	}
	return guiApp.Run()
}

func (g *GUIApp) togglePostOCRFeature(configKey string) {
	cfg := openapi.NewConfiguration()
	cfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client := openapi.NewAPIClient(cfg)

	getResult, _, err := client.ConfigAPI.AppConfigGetGet(context.Background()).Key(configKey).Execute()
	if err != nil {
		g.notifyUser(t("notify.title.error"), fmt.Sprintf(t("error.configUpdateFailed"), err))
		return
	}

	currentValue, ok := getResult["value"].(bool)
	if !ok {
		currentValue = false
	}

	setReq := openapi.NewHandlersAppConfigSetReq(configKey, !currentValue)
	_, _, err = client.ConfigAPI.AppConfigSetPost(context.Background()).HandlersAppConfigSetReq(*setReq).Execute()
	if err != nil {
		g.notifyUser(t("notify.title.error"), fmt.Sprintf(t("error.configUpdateFailed"), err))
		return
	}

	featureName := ""
	switch configKey {
	case "after_ocr.auto_fix_content.enabled":
		featureName = t("ui.postOCR.autoFixLayout")
	case "after_ocr.translate.enabled":
		featureName = t("ui.postOCR.translate")
	case "after_ocr.generate_by_llm.enabled":
		featureName = t("ui.postOCR.generateNotes")
	}

	status := "启用"
	if currentValue {
		status = "禁用"
	}
	g.notifyUser(t("notify.title.config"), fmt.Sprintf("%s 已%s", featureName, status))
}
