# TODO: remove sudo in docker commands
# by delegating the docker build commands to any of the docker user group members

# protobuf specific tasks >> both client code and server interfaces
# the generated code will be used both by the server and client
# hence we populate this into /commons
generate-bindings:
	cd ./commons/proto/ && protoc -I=./ *.proto --go_out=plugins=grpc:../../client/src/go/kaiko/
	cp ./client/src/go/kaiko/*.pb.go ./server/src/go/kaiko/
# client (rpc requester) specific tasks
build-client: | generate-bindings ./client/src/go/main.go
	cd ./client/src/go && go get ./... && go build -o ../../dist/kaiko-exists main.go

run-client: | build-client
	./client/dist/kaiko-exists

# server (rpc responder) specific tasks
build-server: | generate-bindings ./server/src/go/main.go
	cd ./server/src/go && go build -o ../../dist/kaiko-grpc-server main.go

run-server: | build-server
	./client/dist/kaiko-exists

test-server: | build-server
	cd ./server/src/go/kaiko && go test -v

build-server-docker: #| build-server
	cd ./server && sudo docker build -t kaiko-grpc-server .

run-server-docker: #| build-server-docker
	cd ./server && sudo docker run --publish 8080:8080 --name kaiko-grpc-server-poc kaiko-grpc-server

test-server-docker:
	sudo docker exec -it kaiko-grpc-server-poc sh -c "cd /kaiko-grpc-server/src/go/kaiko && go test -v"

kill-server-docker:
	cd ./server && sudo docker kill --signal=SIGTERM kaiko-grpc-server-poc

restart-server-docker:
	cd ./server && sudo docker restart kaiko-grpc-server-poc

remove-server-docker:
	cd ./server && sudo docker rm kaiko-grpc-server-poc

.PHONY: generate-bindings build-client run-client build-server run-server build-docker-server run-server-docker run-server-docker kill-server restart-server remove-server
