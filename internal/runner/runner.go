package runner

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/straydragon/bookxnote-local-ocr/internal/settings"
)

func RunWithDebuggingTLSConfig(r *gin.Engine) {
	// 创建自定义的 TLS 配置
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12, // 只使用 TLS 1.2 ~ 1.3
		MaxVersion:               tls.VersionTLS13,
		InsecureSkipVerify:       true,
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,             // 禁用会话票证
		ClientAuth:               tls.NoClientCert, // 明确不需要客户端证书
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},
		Certificates: make([]tls.Certificate, 1),
		// VerifyConnection: func(cs tls.ConnectionState) error {
		// 	log.Printf("[TLS] 握手完成: Version=%X, CipherSuite=%X", cs.Version, cs.CipherSuite)
		// 	log.Printf("[TLS] 协商的协议版本: %v", cs.NegotiatedProtocol)
		// 	log.Printf("[TLS] 使用的密码套件: %v", tls.CipherSuiteName(cs.CipherSuite))
		// 	log.Printf("[TLS] ServerName: %v", cs.ServerName)
		// 	log.Printf("[TLS] HandshakeComplete: %v", cs.HandshakeComplete)
		// 	return nil
		// },
	}

	// 加载证书链
	cert, err := tls.LoadX509KeyPair(
		settings.GetPathFromCertDir("chain.pem"),
		settings.GetPathFromCertDir("key.pem"),
	)
	if err != nil {
		log.Fatalf("Failed to load certificate chain: %v", err)
	}
	tlsConfig.Certificates[0] = cert

	server := &http.Server{
		Addr:           ":443",
		Handler:        r,
		TLSConfig:      tlsConfig,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		// ConnState: func(conn net.Conn, state http.ConnState) {
		// 	log.Printf("[Conn] %s -> %s: %s",
		// 		conn.RemoteAddr(),
		// 		conn.LocalAddr(),
		// 		state.String())

		// 	if state == http.StateClosed {
		// 		if tlsConn, ok := conn.(*tls.Conn); ok {
		// 			cs := tlsConn.ConnectionState()
		// 			log.Printf("[Conn] 连接关闭详情:")
		// 			log.Printf("  - 版本: %X", cs.Version)
		// 			log.Printf("  - 握手完成: %v", cs.HandshakeComplete)
		// 			log.Printf("  - ServerName: %v", cs.ServerName)
		// 			log.Printf("  - 协商的协议: %v", cs.NegotiatedProtocol)
		// 		}
		// 	}
		// },
	}

	log.Printf("Starting server on :443...")
	if err := server.ListenAndServeTLS(
		settings.GetPathFromCertDir("cert.pem"),
		settings.GetPathFromCertDir("key.pem"),
	); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
