
NAME := cr_connection_checker

build:
	for arch in amd64 arm64; do \
		for os in linux windows; do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o bin/$$os/$$arch/$(NAME) main.go cr_driver.go; \
		done; \
	done

clean:
	rm bin/* -rf

test:
	go test -v ./...

.PHONY: build clean
