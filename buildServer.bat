@echo off
echo Building server ...
cd src/cmd
go build -o ../../build/server.exe
cd ../../
pause