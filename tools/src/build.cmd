@echo off
set SCRIPT=%~dp0
set DIRECTORY=%SCRIPT:~0,-1%
set TARGET_DIRECTORY=%DIRECTORY%\..\bin

set GOOS=windows
set GOARCH=amd64
go build -o "%TARGET_DIRECTORY%\update_hugo.exe" "%DIRECTORY%\update_hugo.go"
if %ERRORLEVEL%==1 goto end

set GOOS=linux
set GOARCH=amd64
go build -o "%TARGET_DIRECTORY%\update_hugo" "%DIRECTORY%\update_hugo.go"
if %ERRORLEVEL%==1 goto end

:end
