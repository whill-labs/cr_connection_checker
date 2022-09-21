
NAME := cr_connection_checker

build:
	for arch in amd64 arm64; do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build -o bin/$$arch/$(NAME) ./cmd/connection_checker; \
	done

clean:
	rm bin/* -rf

.PHONY: build clean

