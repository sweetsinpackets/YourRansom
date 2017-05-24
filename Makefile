all:
	GOOS=windows GOARCH=386 go build -ldflags -H=windowsgui -o=dists/win.exe
	GOOS=windows GOARCH=386 go build -o=dists/win.debug.exe
	GOOS=linux GOARCH=386 go build -o=dists/linux-86
	GOOS=linux GOARCH=amd64 go build -o=dists/linux-64
	GOOS=darwin GOARCH=amd64 go build -o=dists/mac

clean:
	rm -rf dists

upx:
	upx dists/*
