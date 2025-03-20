# 前置准备
- 下载并安装 [UmiOCR](https://github.com/hiroi-sora/Umi-OCR), 解压并运行该程序, 运行该应用保持在后台

# 安装使用
> [!warning]
> 需要完成 [#前置准备](#前置准备)

到 [Release 页面](https://github.com/StrayDragon/bookxnote-local-ocr/releases)下载对应平台的压缩包, 解压后使用命令行/终端运行, 以下命令查看更多帮助
```
./bookxnote-local-ocr -h # 或 Windows 上 .\bookxnote-local-ocr.exe -h
```

以下命令均需要在解压后的程序根目录下运行!

## Linux

1. 下载并解压 tar.gz 文件
2. 进入解压后的目录
3. 运行配置权限脚本：
   ```
   chmod +x setup-privileges.sh
   sudo ./setup-privileges.sh
   ```
4. 安装一些配置, 查看运行提示操作
```sh
./bookxnote-local-ocr install
```

之后仅需要在使用BookxNote OCR功能时, 运行该命令, 保持终端存在, 不要关闭
```sh
./bookxnote-local-ocr server
```

可以使用GUI程序(功能和Cli一致, 但是可以托盘管理)
```sh
./bookxnote-local-ocr-gui
```

## Windows

解压后会看到以下文件：
- `install.bat` - 用于安装证书和配置hosts
- `start.bat` - 打开命令行模式, 每次使用OCR功能时运行, 保持窗口开启
- `start-gui.bat` - 打开图形界面模式, 每次使用OCR功能时运行, 保持窗口开启
- `uninstall.bat` - 清理证书和hosts配置

> [!note]
> - 以上.bat文件都需要管理员权限, 请用管理员权限打开(右键找到以管理员身份运行)

使用步骤：
1. 首次使用时需要运行一次, 使用管理员权限运行`install.bat`, 按提示操作
2. 每次需要使用OCR时, 使用管理员权限运行`start.bat`, 保持窗口开启

也可以使用命令行方式运行, 需要管理员权限
```powershell
.\bookxnote-local-ocr.exe install  # 首次安装
.\bookxnote-local-ocr.exe server   # 启动服务
```

打开 BookXNote, 在右上角选项-文字识别中(需要高级用户)随意配置 API Key 和 Secret Key, 点击应用后检查输出是否为 "应用OCR成功"

如果成功, 则可以正常使用 OCR 功能, 否则参考 [#QA](#QA) 排查问题

# OCR后处理功能(使用GUI(图形界面)模式使用)

本程序支持在OCR识别后对文本进行进一步处理，不过需要提前配置, 包括：

## 自动整理行

OCR识别后的文本可能存在不必要的换行或格式问题，启用此功能可以自动修复文本布局。

## 自动翻译

将OCR识别的文本自动翻译成指定语言，支持多种语言。

## 自动生成笔记

根据OCR识别的文本自动生成结构化笔记，帮助快速整理内容。

## 如何配置

1. 编辑下载程序所在目录下的 `config.yml`文件，配置LLM模型信息：
   ```yaml
   llm:
     models:
       - provider: "openai"
         api_key: "你的OpenAI API密钥"
         api_base_url: "https://api.openai.com/v1"
         name: "gpt-3.5-turbo"
         ident: "openai-gpt35"
   ```

2. 调整或修改 `after_ocr` 中的内容. 如：
   ```yaml
   after_ocr:
     auto_fix_content:
       enabled: true  # 启用自动整理行
     translate:
       enabled: false  # 是否启用自动翻译
       target_language: "zh-CN"  # 目标语言
     generate_by_llm:
       enabled: false  # 是否启用自动生成笔记
       prompt_ident: "summarize"  # 使用的提示模板ID
   ```

3. 使用GUI模式时，可以通过托盘菜单快速开启/关闭这些功能。

# 配置支持本项目标准化的OCR服务接口规范的服务

使用自定义OCR服务, 你需要检查或适配 [openapi spec](../openapi/bookxnote-local-ocr.yaml), 并且运行你的服务

同时需要修改该应用的配置文件如：

```yaml
ocr:
  selected: "custom"
  custom:
    api_base_url: "http://<you-api-base-url>"
    api_key: "your-api-key-here"
  # ...
```

可以参考/使用这个项目提供的server示例 [StrayDragon/llm-ocr-server](https://github.com/StrayDragon/llm-ocr-server/tree/main/llm/got-ocr-2.0-cpu)

# 卸载

在解压后的目录中运行以下命令

## Linux
```
sudo ./bookxnote-local-ocr uninstall
```

## Windows
如果不再使用, 双击运行`uninstall.bat`清理配置 或使用管理员权限打开powershell运行以下命令

```powershell
.\bookxnote-local-ocr.exe uninstall
```
后手动删除安装目录

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

0. 可执行程序所在目录

<!-- ### Linux
1. ~/.local/share/bookxnote-local-ocr/config.yml
2. ~/.config/bookxnote-local-ocr/config.yml

### Windows
1. %APPDATA%/bookxnote-local-ocr/config.yml -->

查看 [config.yml](../config.yml) 或者 [config.yml | Github Repo](https://github.com/StrayDragon/bookxnote-local-ocr/blob/main/artifact/config.yml) 获取更多默认配置信息

## 3. OCR后处理功能不工作?

- 检查配置文件中是否正确配置了LLM模型信息，特别是API密钥
- 检查网络连接是否正常，能否访问OpenAI API
- 检查日志输出，查看是否有错误信息
- 如果使用GUI模式，可以通过托盘菜单重新开启/关闭功能尝试
- 如果仍有问题, 可以[创建issue](https://github.com/StrayDragon/bookxnote-local-ocr/issues/new)
