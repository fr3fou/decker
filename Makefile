PREFIX ?= /usr/local

make:
	make build && make install

build:
	go build cmd/decker/main.go
	mv main decker

install:
	mkdir -p $(PREFIX)/bin
	cp -p decker $(PREFIX)/bin/decker
	chmod 755 $(PREFIX)/bin/decker

uninstall:
	rm -rf $(PREFIX)/bin/decker
