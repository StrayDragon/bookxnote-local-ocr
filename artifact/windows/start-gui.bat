@echo off
echo 正在启动本地OCR服务(请以管理员身份运行), 请在使用时不要关闭该日志窗口, 否则会关闭主程序...
cd /d %~dp0
bookxnote-local-ocr-gui.exe server
pause