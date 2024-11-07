package settings

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetUserConfigDirs 返回用户配置目录列表，按优先级排序
func GetUserConfigDirs() []string {
	var dirs []string

	home, err := os.UserHomeDir()
	if err == nil {
		switch runtime.GOOS {
		case "linux":
			dirs = append(dirs,
				filepath.Join(home, ".config/bookxnote-local-ocr"),
				filepath.Join(home, ".local/share/bookxnote-local-ocr"),
			)
		case "darwin":
			dirs = append(dirs,
				filepath.Join(home, "Library/Application Support/bookxnote-local-ocr"),
			)
		case "windows":
			if appData := os.Getenv("APPDATA"); appData != "" {
				dirs = append(dirs, filepath.Join(appData, "bookxnote-local-ocr"))
			}
		}
	}

	// 添加项目目录作为最后的备选
	dirs = append(dirs, "config")

	return dirs
}

// GetPrimaryUserConfigDir 返回主要的用户配置目录
func GetPrimaryUserConfigDir() string {
	dirs := GetUserConfigDirs()
	if len(dirs) > 0 {
		return dirs[0]
	}
	return "config"
}

// GetCertDir 返回证书存储目录的路径
func GetCertDir() string {
	return filepath.Join(GetPrimaryUserConfigDir(), "certs")
}

// GetPathsFromCertDir 返回证书目录下指定文件的完整路径
func GetPathsFromCertDir(filenames ...string) []string {
	certDir := GetCertDir()
	paths := make([]string, len(filenames))
	for i, filename := range filenames {
		paths[i] = filepath.Join(certDir, filename)
	}
	return paths
}
