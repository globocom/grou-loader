build-linux:
	GOOS=linux GOARCH=amd64 go build -o wrap-loader

build-mac:
	go build -o wrap-loader

clean:
	rm wrap-loader
