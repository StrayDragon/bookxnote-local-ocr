## BookxNotePro 本地OCR服务方案

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/StrayDragon/bookxnote-local-ocr)](./go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/straydragon/bookxnote-local-ocr)](https://goreportcard.com/report/github.com/straydragon/bookxnote-local-ocr)
[![Release](https://img.shields.io/badge/Download-Windows%2FLinux-blue)](https://github.com/StrayDragon/bookxnote-local-ocr/releases)

> [!warning]
> 需要开通 [BookxNote](http://www.bookxnote.com/) 的高级会员才能配置使用

通过劫持[BookxNote](http://www.bookxnote.com/)的OCR请求，直接请求本地OCR服务(默认使用[UmiOCR](https://github.com/hiroi-sora/Umi-OCR))，以获得更好的体验。

![图片](https://github.com/user-attachments/assets/c32b2c54-f678-4865-b7dc-dda607f09787)

同时支持接入用户自定义 OCR API 服务, 只需要实现本项目标准化的OCR服务接口规范([OpenAPI Schema](openapi/ocr-schema.yaml))

除了基础的本地OCR识别, 可以配置 **OpenAI (兼容的)API** 来解锁基于大模型能力的OCR后处理功能!

内置多种实用模板:
- ✨ 自动整理行 - 智能优化OCR文本布局, 让内容更清晰整洁
- 🌍 自动翻译 - 一键将OCR文本翻译成任意语言, 轻松阅读外文资料
- 📝 自动生成笔记 - AI助手帮你总结要点、提炼知识, 事半功倍
- ...

详细使用方式见 [docs/tutorial.md](docs/tutorial.md)

