fmt:
	gofmt -w *.go
test:
	go test -cover

testv:
	go test -cover -v
