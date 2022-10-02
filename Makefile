
NAME := cr_connection_checker

RED=\033[31m
GREEN=\033[32m
RESET=\033[0m
COLORIZE_PASS=sed ''/PASS/s//$$(printf "$(GREEN)PASS$(RESET)")/''
COLORIZE_FAIL=sed ''/FAIL/s//$$(printf "$(RED)FAIL$(RESET)")/''

GO_FILES = $(wildcard *.go)

build:
	for arch in amd64 arm64; do \
		for os in linux windows; do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o bin/$$os/$$arch/$(NAME) $(GO_FILES); \
		done; \
	done

clean:
	rm bin/* -rf

test:
	go test -v ./... | $(COLORIZE_PASS) | $(COLORIZE_FAIL)

.PHONY: build clean
