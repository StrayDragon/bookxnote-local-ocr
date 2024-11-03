.PHONY: all cert clean trust untrust test root_cert

CFG_CERTS_DIR=config/certs
CFG_CERTS_GEN_DIR=config/certs/generated

all: cert trust

root_cert:
	# 1. 生成根证书私钥
	openssl genrsa -out $(CFG_CERTS_GEN_DIR)/rootCA.key 4096

	# 2. 生成根证书
	openssl req -x509 -new -nodes -key $(CFG_CERTS_GEN_DIR)/rootCA.key -sha256 -days 3650 \
		-out $(CFG_CERTS_GEN_DIR)/rootCA.pem \
		-subj "/C=CN/ST=Beijing/L=Beijing/O=LocalDevRootCA/CN=LocalDevRootCA"

cert: $(CFG_CERTS_DIR)/san.cnf root_cert
	# 3. 生成服务器私钥
	openssl genrsa -out $(CFG_CERTS_GEN_DIR)/key.pem 2048

	# 4. 生成证书签名请求(CSR)
	openssl req -new -key $(CFG_CERTS_GEN_DIR)/key.pem -out $(CFG_CERTS_GEN_DIR)/cert.csr -config $(CFG_CERTS_DIR)/san.cnf

	# 5. 使用根证书签发服务器证书
	openssl x509 -req -in $(CFG_CERTS_GEN_DIR)/cert.csr \
		-CA $(CFG_CERTS_GEN_DIR)/rootCA.pem \
		-CAkey $(CFG_CERTS_GEN_DIR)/rootCA.key \
		-CAcreateserial \
		-out $(CFG_CERTS_GEN_DIR)/cert.pem \
		-days 3650 \
		-sha256 \
		-extensions v3_req \
		-extfile $(CFG_CERTS_DIR)/san.cnf

	# 6. 转换为 CRT 格式并创建证书链
	openssl x509 -in $(CFG_CERTS_GEN_DIR)/cert.pem -out $(CFG_CERTS_GEN_DIR)/cert.crt -outform DER
	cat $(CFG_CERTS_GEN_DIR)/cert.pem $(CFG_CERTS_GEN_DIR)/rootCA.pem > $(CFG_CERTS_GEN_DIR)/chain.pem

trust: $(CFG_CERTS_GEN_DIR)/rootCA.pem
	# 信任根证书
	sudo trust anchor --store $(CFG_CERTS_GEN_DIR)/rootCA.pem

untrust:
	sudo rm -f /etc/ca-certificates/trust-source/LocalDevRootCA.p11-kit
	sudo update-ca-trust

test: $(CFG_CERTS_GEN_DIR)/rootCA.pem $(CFG_CERTS_GEN_DIR)/cert.pem
	# 测试证书链
	echo "验证证书链："
	openssl verify -CAfile $(CFG_CERTS_GEN_DIR)/rootCA.pem $(CFG_CERTS_GEN_DIR)/cert.pem

	# 测试 TLS 连接
	echo "\n详细的 TLS 连接测试："
	openssl s_client -connect localhost:443 \
		-servername aip.baidubce.com \
		-CAfile $(CFG_CERTS_GEN_DIR)/rootCA.pem \
		-tls1_2 \
		-debug -msg -state

clean: untrust
	rm -f $(CFG_CERTS_GEN_DIR)/*


run_server:
	@echo "Ensure you have run 'make all' first once."
	go build cmd/server/main.go
	sudo ./main
