# https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
# reducing binary size

build:
	GOOS=windows GOARCH=amd64 GOGC=off \
	go build \
		-ldflags="-s -w" \
		-o ./bin/wap-auth.exe \
		./

# `brew install upx` (mac) or `choco install upx` (win)
pack: build
	upx --brute ./bin/wap-auth.exe