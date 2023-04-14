proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/*.proto

run-agent:
	go run services/agent/*

run-apiserver:
	go run services/apiserver/* config.yaml

run-consul-dev:
	consul agent -server -dev -bind=127.0.0.1 -ui -data-dir=consul-data/

prepare-dev:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/cosmtrek/air@latest
	make proto