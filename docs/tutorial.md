# 编译运行
## Linux

> [!NOTE]
> 仅在 ArchLinux 中测试, 其他发行版未测试

0. 劫持 BookXNote 的 OCR 请求, 在 /etc/hosts 中添加以下内容
```
# BookXNote Pro OCR
127.0.0.1        aip.baidubce.com

```
之后如果不使用了, 记得移除以上配置

1. 安装依赖
- openssl
- go

2. 克隆本仓库, 并运行 `go mod download`

3. 运行 `make all`

- 如果遇到根证书无法安装问题, 请参考 https://wiki.archlinux.org/title/User:Grawity/Adding_a_trusted_CA_certificate 手动信任, 详细见 [makefile](../Makefile)

4. 运行 `make run_server` 启动服务, 需要监听443端口, 因此需要 root 权限

5. 打开 BookXNote, 在右上角选项-文字识别中(需要高级用户)随意配置 API Key 和 Secret Key, 点击应用后检查输出是否为 "应用OCR成功"

6. 下载并安装 [UmiOCR](https://github.com/hiroi-sora/Umi-OCR), 配置 http 服务打开, 运行该应用保持在后台

6. 在 BookXNote 中使用文字识别, 检查是否正常工作


# Q&A

## 1. BookXNote 中配置的 API Key 和 Secret Key 后点击应用, 出现非成功提示?

- 检查是否正确配置了 hosts
- 检查是否已安装并信任自签发的根证书
- 检查 OCR 服务 (如 UmiOCR) 是否正常运行
- 如果以上都没有问题, 请重启 BookXNote , OCR 服务和本服务重试

## 2. 配置文件查找顺序?

程序支持通过配置文件自定义设置。配置文件使用YAML格式, 如

```yaml
ocr:
  umiocr:
    api_url: http://127.0.0.1:1224  # UmiOCR服务地址
```

配置文件位置（按优先级排序）：

### Linux
1. ~/.local/share/bookxnote-local-ocr/config.yml
2. ~/.config/bookxnote-local-ocr/config.yml

### macOS
1. ~/Library/Application Support/bookxnote-local-ocr/config.yml

### Windows
1. %APPDATA%/bookxnote-local-ocr/config.yml

所有平台都支持在项目目录下的 config/config.yml 作为默认配置。