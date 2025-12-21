.PHONY: build install clean run

build:
	go build -o todo

install: build
	cp ./todo ~/.local/bin/todo

clean:
	rm -f todo

run: build
	./todo
