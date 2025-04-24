#!/usr/bin/make -f

APP_NAME = necrolog
CLI_NAME = necrologctl
PKG = ./...

build:
	go build -o $(APP_NAME)
	go build -o $(CLI_NAME) ./necrologctl.go

test:
	go test -v $(PKG)

run:
	go run .

clean:
	rm -f $(APP_NAME) $(CLI_NAME)

install:
	install -m 0755 $(APP_NAME) /usr/local/bin/$(APP_NAME)
	install -m 0755 $(CLI_NAME) /usr/local/bin/$(CLI_NAME)

package:
	fpm -s dir -t deb -n $(APP_NAME) -v 0.1.0 --prefix /usr/local/bin $(APP_NAME) $(CLI_NAME)
