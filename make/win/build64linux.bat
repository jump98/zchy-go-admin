go env -w GOARCH=amd64
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go.exe build -ldflags="-s -w" -o ../bin/radarsvr
set GOOS=windows
set CGO_ENABLED=1
