@echo off
echo 正在启动本地OCR服务(请以管理员身份运行)...
cd /d %~dp0
bookxnote-local-ocr.exe server
pause