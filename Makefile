build:
	go build -o port-forward main.go

test:
	go test -v ./...

clean:
	rm -f port-forward
