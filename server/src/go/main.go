package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"kaiko.io/kaiko" // proto generated + api
)

// var exchangeCode = flag.String("exchange_code", "", "exchange code (eg. 'cbse' for Coinbase or 'bnce' for Binance)")
// var exchangePairCode = flag.String("exchange_pair_code", "", "exchange pair code (eg. 'BTC-USD' for Coinbase or 'BTCUSDT' for Binance)")
// var serverAddress = flag.String("server_address", "", "127.0.0.1:8080")
var port = 8080

func main() {
	// flag.Parse()
	// conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// cli := kaiko.NewKaikoClient(conn)

	// create a listener on TCP localhost port 8080
	println("started listener on port: ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := kaiko.Server{}                       // create a gRPC server object
	grpcServer := grpc.NewServer()            // attach the Ping service to the server
	kaiko.RegisterKaikoServer(grpcServer, &s) // start the server
	println("starting kaiko grpc server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
