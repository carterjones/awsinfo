GO111MODULE := on
unexport GOPATH

build-native:
	go build -mod=vendor -o bin/awsinfo ./cmd/awsinfo/
	go build -mod=vendor -o bin/elbinfo ./cmd/elbinfo/
	go build -mod=vendor -o bin/instanceinfo ./cmd/instanceinfo/
	go build -mod=vendor -o bin/r53info ./cmd/r53info/

build-linux:
	docker run -v $(PWD):/go/src/github.com/carterjones/awsinfo golang:1.11 /bin/bash -c "cd /go/src/github.com/carterjones/awsinfo && make"

clean:
	rm -rf ./bin/

update:
	go get -u
	go mod tidy
	go mod vendor

.PHONY: build-native build-linux clean update
