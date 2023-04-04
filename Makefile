proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/*.proto

run-agent:
	go run services/agent/*

run-apiserver:
	PORT=8080 go run services/apiserver/*