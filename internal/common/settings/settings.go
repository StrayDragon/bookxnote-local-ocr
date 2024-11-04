package settings

import (
	"path/filepath"
)

const (
	CertDir = "config/certs/generated"
)

// GetPathFromCertDir 从配置目录获取证书路径
func GetPathFromCertDir(name string) string {
	return filepath.Join(CertDir, name)
}
