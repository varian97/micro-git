EXECUTABLE := microgit

build:
	go build -o $(EXECUTABLE) main.go
install: build
	sudo mv $(EXECUTABLE) /usr/local/bin/
