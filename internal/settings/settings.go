package settings

import "path/filepath"

const (
	CertDir = "config/certs/generated"
)

func GetPathFromCertDir(name string) string {
	return filepath.Join(CertDir, name)
}
