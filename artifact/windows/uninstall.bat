@echo off
echo 正在卸载配置(请以管理员身份运行)...
cd /d %~dp0
bookxnote-local-ocr.exe uninstall
pause