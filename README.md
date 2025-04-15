## BookxNotePro 本地 OCR 服务方案

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/StrayDragon/bookxnote-local-ocr)](./go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/straydragon/bookxnote-local-ocr)](https://goreportcard.com/report/github.com/straydragon/bookxnote-local-ocr)
[![Release](https://img.shields.io/badge/Download-Windows%2FLinux-blue)](https://github.com/StrayDragon/bookxnote-local-ocr/releases)

> [!warning]
> 需要开通 [BookxNote](http://www.bookxnote.com/) 的高级会员才能配置使用

通过劫持[BookxNote](http://www.bookxnote.com/)的 OCR 请求，直接请求本地 OCR 服务(默认使用[UmiOCR](https://github.com/hiroi-sora/Umi-OCR), 同时支持自建[LLM-OCR-Server](https://github.com/StrayDragon/llm-ocr-server))，以获得更好的体验。

![图片](https://github.com/user-attachments/assets/c32b2c54-f678-4865-b7dc-dda607f09787)

同时支持接入用户自定义 OCR API 服务, 只需要实现本项目标准化的 OCR 服务接口规范([OpenAPI Schema](openapi/bookxnote-local-ocr.yaml))

除了基础的本地 OCR 识别, 可以配置 **OpenAI (兼容的)API** 来解锁基于大模型能力的 OCR 后处理功能!

内置多种实用模板:

- ✨ 自动整理行 - 智能优化 OCR 文本布局, 让内容更清晰整洁
- 🌍 自动翻译 - 一键将 OCR 文本翻译成任意语言, 轻松阅读外文资料
- 📝 自动生成笔记 - AI 助手帮你总结要点、提炼知识, 事半功倍
- ...

详细使用方式见 [docs/tutorial.md](docs/tutorial.md)
