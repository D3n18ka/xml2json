.PHONY : clean server docker

all: clean server

clean:
	rm -f server

server:
	go build -o server ./cmd/server/

docker:
	docker build --ssh default -t server:latest .
