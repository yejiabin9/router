
.PHONY: proto
proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) cap1573/cap-v3 --proto_path=. --micro_out=. --go_out=:. ./proto/router/router.proto

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o router *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t router:latest
