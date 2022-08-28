
cd ..
go mod tidy

Write-Output "正在编译 windows 版本"
go build -ldflags="-s -w" -o build/smq.exe

Write-Output "正在编译 mac 版本"
go env -w GOOS=darwin
go build -ldflags="-s -w" -o build/smq-darwin

Write-Output "正在编译 linux 版本"
go env -w GOOS=linux
go build -ldflags="-s -w" -o build/smq-linux

cd build
upx .\smq.exe
upx .\smq-darwin
upx .\smq-linux

go env -w GOOS=windows
