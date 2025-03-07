//go:generate fyne bundle -o gui_resources.go ../../artifact/icon.png
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

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

		"notify.title.server":       "服务器状态",
		"notify.title.error":        "错误",
		"notify.title.installing":   "安装中",
		"notify.title.success":      "成功",
		"notify.title.uninstalling": "卸载中",

		"notify.msg.serverStarted":     "服务器已启动",
		"notify.msg.serverStopped":     "服务器已停止",
		"notify.msg.installing":        "正在安装配置...",
		"notify.msg.installComplete":   "配置安装完成",
		"notify.msg.uninstalling":      "正在卸载配置...",
		"notify.msg.uninstallComplete": "配置卸载完成",

		"error.serverAlreadyRunning": "服务器已在运行",
		"error.serverNotRunning":     "服务器未运行",
		"error.startFailed":          "启动服务器失败: %v",
		"error.stopFailed":           "停止服务器失败: %v",
		"error.installFailed":        "安装配置失败: %v",
		"error.uninstallFailed":      "卸载配置失败: %v",

		"log.serverStarted":  "[bookxnote-local-ocr] %v | 服务器已启动\n",
		"log.serverStopped":  "[bookxnote-local-ocr] %v | 服务器已停止\n",
		"log.serverStopping": "[bookxnote-local-ocr] %v | 已停止服务器...\n",
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
	*App
	fyneApp       fyne.App
	window        fyne.Window
	serverProcess *os.Process
	serverSwitch  *widget.Check
}

func NewGUIApp() (*GUIApp, error) {
	cliApp, err := NewApp()
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

func (g *GUIApp) handleServerStart(statusLabel *widget.Label) {
	if g.serverProcess != nil {
		return
	}

	go func() {
		statusLabel.SetText(t("status.starting"))

		cmd := exec.Command(g.executablePath, "server")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			dialog.ShowError(fmt.Errorf(t("error.startFailed"), err), g.window)
			g.serverSwitch.SetChecked(false)
			statusLabel.SetText(t("status.startFailed"))
			return
		}

		g.serverProcess = cmd.Process
		statusLabel.SetText(t("status.running"))
		g.notifyUser(t("notify.title.server"), t("notify.msg.serverStarted"))
		fmt.Printf(t("log.serverStarted"), time.Now().Format("2006/01/02 - 15:04:05"))

		go func() {
			_ = cmd.Wait()
			g.serverProcess = nil
			g.serverSwitch.SetChecked(false)

			statusLabel.SetText(t("status.stopped"))
			g.notifyUser(t("notify.title.server"), t("notify.msg.serverStopped"))
			fmt.Printf(t("log.serverStopped"), time.Now().Format("2006/01/02 - 15:04:05"))
		}()
	}()
}

func (g *GUIApp) handleServerStop(statusLabel *widget.Label) {
	if g.serverProcess == nil {
		return
	}

	statusLabel.SetText(t("status.stopping"))
	if err := g.serverProcess.Signal(os.Interrupt); err != nil {
		dialog.ShowError(fmt.Errorf(t("error.stopFailed"), err), g.window)
		g.serverSwitch.SetChecked(true)
		return
	}

	statusLabel.SetText(t("status.stopped"))
}

func (g *GUIApp) handleInstallConfig() {
	dialog.ShowInformation(t("notify.title.installing"), t("notify.msg.installing"), g.window)

	cmd := exec.Command(g.executablePath, "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		dialog.ShowError(fmt.Errorf(t("error.installFailed"), err), g.window)
	} else {
		dialog.ShowInformation(t("notify.title.success"), t("notify.msg.installComplete"), g.window)
	}
}

func (g *GUIApp) handleUninstallConfig() {
	dialog.ShowInformation(t("notify.title.uninstalling"), t("notify.msg.uninstalling"), g.window)

	cmd := exec.Command(g.executablePath, "uninstall")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
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
		return fmt.Errorf(t("error.serverAlreadyRunning"))
	}

	cmd := exec.Command(g.executablePath, "server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf(t("error.startFailed"), err)
	}

	g.serverProcess = cmd.Process
	g.notifyUser(t("notify.title.server"), t("notify.msg.serverStarted"))
	fmt.Printf(t("log.serverStarted"), time.Now().Format("2006/01/02 - 15:04:05"))

	if g.serverSwitch != nil {
		g.serverSwitch.SetChecked(true)
	}

	go g.monitorServerProcess(cmd)

	return nil
}

func (g *GUIApp) monitorServerProcess(cmd *exec.Cmd) {
	_ = cmd.Wait()
	g.serverProcess = nil

	if g.serverSwitch != nil {
		g.serverSwitch.SetChecked(false)
	}

	fmt.Printf(t("log.serverStopped"), time.Now().Format("2006/01/02 - 15:04:05"))
}

func (g *GUIApp) stopServer() error {
	if g.serverProcess == nil {
		return fmt.Errorf(t("error.serverNotRunning"))
	}

	if err := g.serverProcess.Signal(os.Interrupt); err != nil {
		return fmt.Errorf(t("error.stopFailed"), err)
	}

	g.notifyUser(t("notify.title.server"), t("notify.msg.serverStopped"))
	fmt.Printf(t("log.serverStopping"), time.Now().Format("2006/01/02 - 15:04:05"))

	tempProcess := g.serverProcess
	g.serverProcess = nil

	if g.serverSwitch != nil {
		g.serverSwitch.SetChecked(false)
	}

	go g.verifyServerStopped(tempProcess)

	return nil
}

func (g *GUIApp) verifyServerStopped(process *os.Process) {
	time.Sleep(500 * time.Millisecond)

	if process == nil {
		return
	}

	// Signal(0) doesn't send a signal but checks if process exists
	if process.Signal(syscall.Signal(0)) == nil {
		g.serverProcess = process

		if g.serverSwitch != nil {
			g.serverSwitch.SetChecked(true)
		}
	}
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
		desktopApp.SetSystemTrayMenu(createTrayMenu(g, toggleServerItem))
	}

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
	return fyne.NewMenu(t("appTitle"),
		toggleServerItem,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(t("ui.show"), func() {
			g.window.Show()
		}),
		fyne.NewMenuItem(t("ui.exit"), func() {
			if g.serverProcess != nil {
				_ = g.stopServer()
			}
			g.fyneApp.Quit()
		}),
	)
}

func (g *GUIApp) updateMenuPeriodically(desktopApp desktop.App, toggleServerItem *fyne.MenuItem, updateMenu func()) {
	var lastServerState bool
	for {
		time.Sleep(1 * time.Second)
		if g.fyneApp == nil {
			return
		}

		// Only update the menu if the server state has changed
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

func runGUI() error {
	guiApp, err := NewGUIApp()
	if err != nil {
		return err
	}
	return guiApp.Run()
}
