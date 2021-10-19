@echo Building...
@echo off

set RSRC=%GOPATH%\bin\rsrc.exe
if not exist %RSRC% (
go get -u github.com/akavel/rsrc
)

set GOOS=windows
::set GOARCH=386

%RSRC% -manifest="nac.manifest" -ico="syncthing.ico" -o="rsrc.syso"

go build -ldflags "-s -w" -o autohosts.exe
del /F /Q rsrc.syso >nul 2>nul
