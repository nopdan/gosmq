
cd ..
xcopy web\dist\ internal\serve\dist\ /E /Y

go mod tidy

Write-Output "编译 windows 版本"
go build -o build/smq.exe

cd build
