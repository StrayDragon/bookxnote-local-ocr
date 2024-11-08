package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/straydragon/bookxnote-local-ocr/internal/cert"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
)

func certificatesExist(certDir string) bool {
	files := []string{"rootCA.pem", "rootCA.key", "cert.pem", "key.pem"}
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(certDir, file)); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func loadRootCertFromFile(certPath string) (*x509.Certificate, error) {
	pemBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	return x509.ParseCertificate(block.Bytes)
}

func main() {
	forceRegen := flag.Bool("force", false, "强制重新生成证书")
	clean := flag.Bool("clean", false, "清理已生成的证书")
	flag.Parse()

	certDir := settings.GetCertDir()
	if err := os.MkdirAll(certDir, 0755); err != nil {
		log.Fatalf("创建证书目录失败 > %v", err)
	}

	if *forceRegen && *clean {
		log.Fatal("无法同时使用 -force 和 -clean 参数")
	}

	if *forceRegen {
		*clean = true
	}

	if *clean {
		rootCertPath := filepath.Join(certDir, "rootCA.pem")
		if _, err := os.Stat(rootCertPath); err == nil {
			if rootCert, err := loadRootCertFromFile(rootCertPath); err == nil {
				if err := cert.UninstallRootCert(rootCert); err != nil {
					log.Printf("卸载根证书失败 > %v", err)
				} else {
					log.Println("已卸载根证书")
				}
			}
		}

		files := []string{"rootCA.pem", "rootCA.key", "cert.pem", "key.pem"}
		for _, file := range files {
			path := filepath.Join(certDir, file)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				log.Printf("清理证书文件失败 %s > %v", path, err)
			}
		}
		log.Println("证书文件已清理")
		if !*forceRegen {
			return
		}
	}

	if certificatesExist(certDir) && !*forceRegen {
		log.Println("证书文件已存在，跳过生成。使用 -force 参数强制重新生成")
		return
	}

	config := cert.CertConfig{
		Country:      []string{"CN"},
		Province:     []string{"Beijing"},
		Locality:     []string{"Beijing"},
		Organization: []string{"Local OCR"},
		CommonName:   "aip.baidubce.com",
		DNSNames: []string{
			"aip.baidubce.com",
			"*.baidubce.com",
			"localhost",
		},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
		CertDir: certDir,
	}

	// 生成根证书+服务器证书
	rootCert, rootKey, err := cert.GenerateRootCA(config)
	if err != nil {
		log.Fatalf("生成根证书失败 > %v", err)
	}
	serverCert, serverKey, err := cert.GenerateServerCert(config, rootCert, rootKey)
	if err != nil {
		log.Fatalf("生成服务器证书失败 > %v", err)
	}

	// 保存到用户配置中一份儿 方便管理
	if err := cert.SaveCertAndKey(rootCert, rootKey,
		filepath.Join(certDir, "rootCA.pem"),
		filepath.Join(certDir, "rootCA.key")); err != nil {
		log.Fatalf("保存根证书失败 > %v", err)
	}
	if err := cert.SaveCertAndKey(serverCert, serverKey,
		filepath.Join(certDir, "cert.pem"),
		filepath.Join(certDir, "key.pem")); err != nil {
		log.Fatalf("保存服务器证书失败 > %v", err)
	}

	// 安装并信任根证书
	if err := cert.InstallRootCert(rootCert); err != nil {
		log.Fatalf("安装根证书失败 > %v", err)
	}

	log.Println("证书生成和安装完成")
}
