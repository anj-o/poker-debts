.PHONY: all run

all: run

run:
	templ generate
	go run .