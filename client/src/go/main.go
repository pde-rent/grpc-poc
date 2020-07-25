package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"kaiko.io/kaiko"
)

var exchangeCode = flag.String("exchange_code", "", "exchange code (eg. 'cbse' for Coinbase or 'bnce' for Binance)")
var exchangePairCode = flag.String("exchange_pair_code", "", "exchange pair code (eg. 'BTC-USD' for Coinbase or 'BTCUSDT' for Binance)")
var serverAddress = flag.String("server_address", "", "address for the kaiko server (eg. 127.0.0.1:8080)")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	cli := kaiko.NewKaikoClient(conn)

	res, err := cli.Exists(context.Background(), &kaiko.ExistsRequest{
		ExchangeCode:     *exchangeCode,
		ExchangePairCode: *exchangePairCode,
	})
	if err != nil {
		log.Fatal(err)
	}

	switch res.Exists {
	case kaiko.ExistsResponse_YES:
		println("OK")
	case kaiko.ExistsResponse_NO:
		println("KO")
	default:
		log.Fatal("unexpected response")
	}
}
