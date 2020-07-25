package kaiko

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

var testKaikoClient = defaultTestClient()
var existingExchanges = getExistingExchanges()
var existingInstrumentsByExchanges = getExistingInstrumentsByExchanges()

// TODO: replace the fmt.Printf and printl with proper logs in logfiles

// default integration client
func defaultTestClient() KaikoClient {
	conn, err := grpc.Dial(GRPC_URL, grpc.WithInsecure())
	if err != nil {
		println(err)
	}
	cli := NewKaikoClient(conn)
	return cli
}

// generate random instrument name along a valid exchange
func getFakeInstrument() (string, string) {

	return getExistingExchange(), string(saltGenerate(16))
}

// generate random instrument name
func getExistingExchange() string {

	return existingExchanges[rand.Intn(len(existingExchanges))]
}

// generate random valid exchange along vali dinstrument
func getExistingInstrument() (string, string) {

	i := rand.Intn(len(existingInstrumentsByExchanges))
	for k, v := range existingInstrumentsByExchanges {
		if i == 0 {
			return k, v[rand.Intn(len(v))]
		}
		i--
	}
	return "", ""
}

// get list of all valid exchanges
func getExistingExchanges() []string {

	exchanges, err := getExchanges()
	if err != nil {
		return nil
	}
	return exchanges
}

// get list of all valid exchanges with their instruments
func getExistingInstrumentsByExchanges() map[string][]string {

	err := refreshInstrumentsByExchange()
	if err != nil {
		return nil
	}
	return instrumentsByExchange
}

// TODO: implement batch + loading from test files
func testExistsSingle(exchangeCode string, exchangePairCode string) (bool, int, error) {

	// time at request
	start := time.Now().UnixNano()
	// execute gRPC call on given function
	res, err := testKaikoClient.Exists(context.Background(), &ExistsRequest{
		ExchangeCode:     exchangeCode,
		ExchangePairCode: exchangePairCode,
	})
	// time at response
	end := time.Now().UnixNano()
	elapsed := int((end - start) / 1000000)
	// check for a call bad return
	if err != nil {
		println(err)
	}
	// analyse result
	switch res.Exists {
	// instrument exists for exchange
	case ExistsResponse_YES:
		return true, elapsed, nil
	// instrument does not exist for exchange
	case ExistsResponse_NO:
		return false, elapsed, nil
	// bad RPC call or response (out of the proto interface)
	default:
		// disregard the 2 firsts return values
		return false, 0, errors.New("unexpected response")
	}
}
