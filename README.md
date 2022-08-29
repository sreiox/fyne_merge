## 




### 打包

> macos 打包 Windows 需要执行安装 `brew install mingw-w64` [在macOS下启用CGO_ENABLED的交叉编译Go语言项目生成Windows EXE文件](https://blog.csdn.net/mctlilac/article/details/105605147)


#### 在macos arm 架构下打包命令如下

```shell
## 打包 Windows 可执行文件
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w"

## 打包 macos arm 软件
go build main.go

## 打包 macos amd 软件

```