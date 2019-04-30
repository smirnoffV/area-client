install:
	go mod download
	git submodule update --remote --init

build-proto:
	rm -rf ./pb/*
	protoc -I protocol protocol/*.proto --go_out=plugins=grpc:./pb

build: build-proto
	rm -rf ./bin/*
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o=./bin/linux-app-64 ./cmd/main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o=./bin/darwin-app-64 ./cmd/main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o=./bin/win-app-64.exe ./cmd/main.go

docker-build: build
	docker build -t area-svc .
	docker tag area-svc smirnoffv/area-svc:dev

docker-publish:
	docker push smirnoffv/area-svc:dev