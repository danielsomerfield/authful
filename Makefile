OSX_BINARY=authful-osx
LINUX_BINARY=authful-linux-386

build: test
	go build -o ${OSX_BINARY}
	GOOS=linux GOARCH=386 go build -o ${LINUX_BINARY}

install_deps:
	go get golang.org/x/oauth2
	go get golang.org/x/crypto/scrypt
	go get gopkg.in/yaml.v2
	go get github.com/PuerkitoBio/goquery

test: test-only

test-only:
	go test ./...

clean:
	go clean -r -i
	rm -f ${OSX_BINARY} ${LINUX_BINARY}


