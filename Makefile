.PHONY: all

all: debug release

debug:
	go build -o target/dbg.exe

release:
	go build -ldflags "-s -w" -o target/minall.exe