cd ../src

GOOS=windows GOARCH=amd64 go build -o ../build/wsgo_windwos_x64.exe

GOOS=windows GOARCH=386 go build -o ../build/wsgo_windwos_x86.exe

GOOS=linux GOARCH=amd64 go build -o ../build/wsgo_linux_x64

GOOS=linux GOARCH=386 go build -o ../build/wsgo_linux_x86

GOOS=darwin GOARCH=amd64 go build -o ../build/wsgo_mac_x64

GOOS=darwin GOARCH=386 go build -o ../build/wsgo_mac_x86
