@echo off
cd ../src
set GOOS=windows
set GOARCH=amd64
echo "build window x64 version"
go build -o ../build/wsgo_windows_x64.exe
echo "build finish wsgo_windows_x64.exe"

set GOOS=linux
set GOARCH=amd64
echo "build linux x64 version"
go build -o ../build/wsgo_linux_x64
echo "build finish wsgo_linux_x64"

set GOOS=darwin
set GOARCH=amd64
echo "build Mac x64 version"
go build -o ../build/wsgo_mac_x64
echo "build finish wsgo_mac_x64"

set GOOS=windows
set GOARCH=386
echo "build window x86 version"
go build -o ../build/wsgo_windows_x86.exe
echo "build finish wsgo_windows_x86.exe"

set GOOS=linux
set GOARCH=386
echo "build linux x86 version"
go build -o ../build/wsgo_linux_x86
echo "build finish wsgo_linux_x86"

set GOOS=darwin
set GOARCH=386
echo "build Mac x86 version"
go build -o ../build/wsgo_mac_x86
echo "build finish wsgo_mac_x86"
