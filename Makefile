build: ./main.go
	go build -o bin/cloudctl main.go

test: pkg/*/*.go cmd/*.go
	go test ./...

vtest: pkg/*/*.go cmd/*.go
	go test -v ./...
# go test -v github.com/cloudor-io/cloudctl/pkg/...

clean:
	rm -rf bin/*

module:
	go mod vendor