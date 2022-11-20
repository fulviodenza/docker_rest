run:
	touch ./tmp/interrupt_task.txt
	go run main.go ubuntu:latest

test:
	cd ./internal/docker_client && go test