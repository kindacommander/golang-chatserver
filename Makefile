PKG := github.com/kindacommander/golang-chatserver

go-build:
	@echo " > Building binary..."
	go build -o vendor/chat cmd/chat/main.go

BIN=$(shell pwd)/vendor/chat
exec:
	@echo " > Executing..."
ifneq ("$(wildcard $(BIN))","")
	vendor/chat
else
	@echo " Error: No binary file"
endif