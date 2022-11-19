run:
	touch ./tmp/interrupt_task.txt
	go run main.go

test:
	cd ./internal/docker_client && go test