package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/smallstep/truststore"
)

type CertConfig struct {
	Country      []string
	Province     []string
	Locality     []string
	Organization []string
	CommonName   string
	DNSNames     []string
	IPAddresses  []net.IP
	CertDir      string
}

// GenerateRootCA 生成根证书
func GenerateRootCA(config CertConfig) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("生成私钥失败 > %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      config.Country,
			Province:     config.Province,
			Locality:     config.Locality,
			Organization: []string{"LocalDevRootCA"},
			CommonName:   "LocalDevRootCA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 使用生成的私钥签名
	derBytes, err := x509.CreateCertificate(rand.Reader,
		template,
		template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("生成根证书失败 > %v", err) // 自签名根证书
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("解析证书失败 > %v", err)
	}

	return cert, privateKey, nil
}

// GenerateServerCert 生成服务器证书, 并使用rootCA签名
func GenerateServerCert(config CertConfig, rootCA *x509.Certificate, rootKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("生成私钥失败 > %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Country:      config.Country,
			Province:     config.Province,
			Locality:     config.Locality,
			Organization: config.Organization,
			CommonName:   config.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              config.DNSNames,
		IPAddresses:           config.IPAddresses,
	}

	// 使用根证书签名
	derBytes, err := x509.CreateCertificate(rand.Reader,
		template,
		rootCA,
		&privateKey.PublicKey,
		rootKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("生成服务器证书失败 > %v", err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("解析证书失败 > %v", err)
	}

	return cert, privateKey, nil
}

// SaveCertAndKey 保存证书和私钥到文件
func SaveCertAndKey(cert *x509.Certificate, key *rsa.PrivateKey, certPath, keyPath string) error {
	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		return fmt.Errorf("创建证书目录失败 > %v", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("创建证书文件失败 > %v", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
		return fmt.Errorf("写入证书失败 > %v", err)
	}

	keyOut, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败 > %v", err)
	}
	defer keyOut.Close()

	if err := pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return fmt.Errorf("写入私钥失败 > %v", err)
	}

	return nil
}

// InstallRootCert 安装并信任根证书
func InstallRootCert(rootCert *x509.Certificate) error {
	err := truststore.Install(rootCert)
	if err != nil {
		return fmt.Errorf("安装根证书失败 > %v", err)
	}
	return nil
}

// UninstallRootCert 卸载根证书
func UninstallRootCert(rootCert *x509.Certificate) error {
	err := truststore.Uninstall(rootCert)
	if err != nil {
		return fmt.Errorf("卸载根证书失败, 请尝试自己卸载, 原因: %v", err)
	}
	return nil
}
