proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/*.proto

run-agent:
	go run services/agent/*

run-apiserver:
	PORT=8080 go run services/apiserver/*

run-consul-dev:
	consul agent -server -dev -bind=127.0.0.1 -ui -data-dir=consul-data/