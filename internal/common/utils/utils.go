package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GetExecLinuxCaps (Linux) 获取当前进程的 Linux capabilities
func GetExecLinuxCaps(pid int) (string, error) {
	// $ getpcaps <PID>
	// PID: = cap_chown,cap_dac_override,cap_fowner,cap_fsetid,cap_kill,cap_setgid,cap_setuid,cap_setpcap,cap_net_raw,cap_sys_chroot,cap_mknod,cap_audit_write,cap_setfcap+i
	result, err := exec.Command("getpcaps", fmt.Sprintf("%d", pid)).Output()
	if err != nil {
		return "", err
	}

	return string(bytes.ToLower(result)), nil
}

// CheckCurrentProcessCaps (Linux) 检查当前进程是否具有所需的 Linux capabilities
func CheckCurrentProcessCaps(requiredCaps []string) error {
	pid := os.Getpid()

	capResult, err := GetExecLinuxCaps(pid)
	if err != nil {
		return fmt.Errorf("检查权限失败: %w", err)
	}

	missingCaps := []string{}
	for _, cap := range requiredCaps {
		if !strings.Contains(capResult, strings.ToLower(cap)) {
			missingCaps = append(missingCaps, cap)
		}
	}

	if len(missingCaps) > 0 && !IsRunningAsRoot() {
		return fmt.Errorf("缺少所需权限: %v，请使用 sudo 运行或设置相应的 Linux capabilities", missingCaps)
	}

	return nil
}

// IsRunningAsRoot (Linux) 检查当前进程是否以 root 用户运行
func IsRunningAsRoot() bool {
	uid := os.Getuid()
	return uid == 0
}

// CheckAdminPrivileges 检查当前进程是否具有管理员权限
func CheckAdminPrivileges() error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("net", "session")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("需要管理员权限运行! 请右键点击程序，选择「以管理员身份运行」")
		}
	} else if runtime.GOOS == "linux" {
		// if !IsRunningAsRoot() {
		// 	return fmt.Errorf("需要管理员权限运行! 请使用 sudo 运行此程序")
		// }
		return nil
	}
	return nil
}

// AllValueOfMap 检查 map 中的所有值是否都满足给定的条件, NOTE: 如果 map 为空, 则返回 true, 否则如果 map 中的值都不满足给定的条件, 则返回 false
func AllValueOfMap[K comparable, V any](m map[K]V, f func(V) bool) bool {
	for _, v := range m {
		if !f(v) {
			return false
		}
	}
	return true
}

// GetExecDir 获取当前可执行程序运行目录
func GetExecDir() (string, error) {
	execDir, err := os.Executable()
	if err == nil {
		return filepath.Dir(execDir), nil
	}
	return "", err
}
