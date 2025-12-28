.PHONY: build install clean run quick

build:
	go build -o todo

install: build
	cp ./todo ~/.local/bin/todo

clean:
	rm -f todo

run: build
	./todo

quick: build
	./todo --quick
