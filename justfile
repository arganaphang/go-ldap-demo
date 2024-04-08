set dotenv-load

# Format Golang
format:
	gofumpt -l -w .
	goimports-reviser -rm-unused -set-alias ./...
	golines -w -m 120 *.go

# build -> build application
build:
	go build -o ./.tmp/main ./cmd/server

# run -> application
run:
	./.tmp/main

# dev -> run build then run it
dev: 
	watchexec -r -c -e go -- just build run

# health -> Hit Health Check Endpoint
health:
	curl -s http://localhost:8000/healthz | jq
