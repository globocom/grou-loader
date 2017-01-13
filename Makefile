build-linux:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o wrap-loader

build-mac:
	CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o wrap-loader

clean:
	rm wrap-loader
