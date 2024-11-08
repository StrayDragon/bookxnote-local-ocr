@echo off
echo 正在启动本地OCR服务...
powershell -Command "Start-Process -Verb RunAs cmd -ArgumentList '/c cd /d %~dp0 && bookxnote-local-ocr.exe server'"
pause