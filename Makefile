build-cli-mac:
	env GOOS=darwin  go build -ldflags="-s -w" -o bin/mac/timesh cmd/cli/main.go
build-cli-win:
	env GOOS=windows go build -ldflags="-s -w" -o bin/win/timesh.exe cmd/cli/main.go
