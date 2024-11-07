# 前置准备
- 配置 hosts 文件, 劫持 BookXNote 的 OCR 请求, 添加以下内容到, **之后如果不使用了, 记得移除以上配置, 并调用相关清理卸载命令, 以免影响正常使用**
  - Linux/macOS: /etc/hosts
  - Windows: C:\Windows\System32\drivers\etc\hosts
  ```
  # BookXNote Pro OCR
  127.0.0.1        aip.baidubce.com
  ```
- 下载并安装 [UmiOCR](https://github.com/hiroi-sora/Umi-OCR), 配置 http 服务打开, 运行该应用保持在后台

# 安装使用
> [!warning]
> 需要完成 [#前置准备](#前置准备)

到 Release 页面下载对应平台的压缩包, 解压后使用命令行/终端运行

## Linux/macOS
```sh
./bookxnote-local-ocr server
```

## Windows
```powershell
.\bookxnote-local-ocr.exe server
```


将自动配置证书并启动服务器, 以下命令查看更多帮助
```
bookxnote-local-ocr -h
```

打开 BookXNote, 在右上角选项-文字识别中(需要高级用户)随意配置 API Key 和 Secret Key, 点击应用后检查输出是否为 "应用OCR成功"

如果成功, 则可以正常使用 OCR 功能, 否则参考 [#QA](#QA) 排查问题

# 本地开发运行
## Linux
> [!warning]
> 需要完成 [#前置准备](#前置准备)

- 安装并正确配置 [Golang](https://go.dev/doc/install)
- 运行 `make build` 会自动下载依赖并完成编译
- 浏览代码, 进行开发 :)


# QA

## 1. BookXNote 中配置的 API Key 和 Secret Key 后点击应用, 出现非成功提示?

- 检查是否正确配置了 hosts
- 检查是否已安装并信任自签发的根证书
- 检查 OCR 服务 (如 UmiOCR) 是否正常运行
- 如果以上都没有问题, 请重启 BookXNote , OCR 服务和本服务重试

ArchLinux 用户如果遇到根证书无法安装问题, 请参考 https://wiki.archlinux.org/title/User:Grawity/Adding_a_trusted_CA_certificate 手动信任, 详细见 [makefile](../Makefile)


## 2. 配置文件查找顺序?

程序支持通过配置文件自定义设置。配置文件使用YAML格式, 如

```yaml
# ...
ocr:
  umiocr:
    api_url: http://127.0.0.1:1224  # UmiOCR服务地址
# ...
```

配置文件位置（按优先级排序）：

### Linux
1. ~/.local/share/bookxnote-local-ocr/config.yml
2. ~/.config/bookxnote-local-ocr/config.yml

### macOS
1. ~/Library/Application Support/bookxnote-local-ocr/config.yml

### Windows
1. %APPDATA%/bookxnote-local-ocr/config.yml

查看 [config/config.yml](../config/config.yml) 获取更多默认配置信息