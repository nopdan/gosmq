
cd frontend
bun install
bun run build
cd ..
rm pkg/server/dist -r
xcopy frontend\dist\ pkg\server\dist\ /E /Y

go mod tidy

Write-Output "编译 windows 版本"
go build -ldflags="-s -w" -o build/smq.exe
