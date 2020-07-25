# Requirements

* Make
* [Go 1.13.x](https://golang.org/dl/)
* [Protobuf 1.3.x](https://github.com/gogo/protobuf/releases)
* [Go protobuf](https://github.com/golang/protobuf) `go get -u github.com/golang/protobuf/protoc-gen-go`

## Build

Once the requirements are installed, just run `make build-client`, it will output a `kaiko-exists` binary for your platform.

## How to use the client

* `$ ./kaiko-exists -exchange_code cbse -exchange_pair_code btc-usd -server_address 127.0.0.1:8080` should output "OK"
* `$ ./kaiko-exists -exchange_code cbse -exchange_pair_code foo-bar` should output "KO"
* `$ ./kaiko-exists -exchange_code foobar -exchange_pair_code btc-usd` should exit with an error code and an "unexpected response" message 
