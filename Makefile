ENV := dev

include .env
export

.PHONY: start
start:
	go run cmd/api/server.go --env=${ENV}

.DEFAULT_GOAL := start
