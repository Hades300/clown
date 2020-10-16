buildall:
	GOOS=windows GOARCH=amd64 go build -o clown_windows_amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o clown_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -o clown_darwin_amd64
build:
	go build .