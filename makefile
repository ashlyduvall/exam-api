#
# Makefile for handling local builds, runs and soforth.
#
# @author Ashly Duvall
#

edit:
	@export GOPATH=$(shell pwd) && cd src/main && "$${EDITOR:-vim}" .

install:
	@export GOPATH=$(shell pwd) && cd src/main && go get .

build:
	@export GOPATH=$(shell pwd) && cd src/main && \
		go build -o inventory . && \
		mv inventory ../../

run:
	@export GOPATH=$(shell pwd) && cd src/main && \
		DB_HOST=$(DB_HOST) \
		DB_PORT=$(DB_PORT) \
		DB_USER=$(DB_USER) \
		DB_PASS=$(DB_PASS) \
		DB_SCHEMA=$(DB_SCHEMA) \
		go run .
