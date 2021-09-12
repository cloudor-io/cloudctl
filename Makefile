build: ./main.go
	GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/cloudor main.go

release: ./main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/cloudor main.go

test: pkg/*/*.go cmd/*.go
	go test ./...

vtest: pkg/*/*.go cmd/*.go
	go test -v ./...
# go test -v github.com/cloudor-io/cloudctl/pkg/...

clean:
	rm -rf bin/*

module:
	go mod vendor