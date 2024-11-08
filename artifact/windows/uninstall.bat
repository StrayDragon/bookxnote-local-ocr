@echo off
echo 正在卸载配置...
powershell -Command "Start-Process -Verb RunAs cmd -ArgumentList '/c cd /d %~dp0 && bookxnote-local-ocr.exe uninstall'"
pause