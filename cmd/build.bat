set GOOS = windows
set GOARCH = 386
go build -ldflags "-s -w -H windowsgui" -trimpath -o "Terraria-Editor.exe"

