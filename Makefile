.PHONY: build

build: linux osx

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/linux/acs main.go
	cd dist/linux && tar cvfz acs-linux.tar.gz acs

osx:
	GOOS=darwin GOARCH=amd64 go build -o dist/osx/acs main.go
	cd dist/osx && tar cvfz acs-osx.tar.gz acs
