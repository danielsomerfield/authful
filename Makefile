build: test
	go build -o authful

install_deps:
	go get golang.org/x/oauth2

test: install_deps
	go test ./...

clean:
	go clean -r -i

