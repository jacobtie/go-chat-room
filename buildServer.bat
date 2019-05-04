@echo off
echo Building server ...
cd backend/cmd
go build -o ../../build/server.exe
cd ../../
pause