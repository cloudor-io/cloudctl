build: ./main.go
	go build -o bin/cloudctl main.go

test: pkg/*/*.go
	go test -v cloudctl/pkg/...

clean:
	rm -rf bin/*

module:
	go mod vendor