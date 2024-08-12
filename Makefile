EXECUTABLE := microgit

build:
	go build -o $(EXECUTABLE) cmd/cmd.go
install: build
	sudo mv $(EXECUTABLE) /usr/local/bin/
