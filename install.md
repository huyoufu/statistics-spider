### windows

```bat
##编译为linux
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go  build -o ss *.go
##编译为macosx
set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=amd64
go build -o ss github.com/huyoufu/statistics-spider
##编译为windows
SET GOOS=windows
go build -o ss.exe  *.go
```

### macosx

```shell
##编译为windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ss.exe main.go
##编译为linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ss main.go
```

### linux

```shell
##编译为windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ss.exe main.go
##编译为macosx
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ss main.go
```
**Enjoy!**