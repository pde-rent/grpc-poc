package kaiko

import (
	"encoding/json"
	"errors"
	fmt "fmt"
	"net/http"
	"time"
)

// global static ionstrument cache, not to overload the rest api
var instrumentsCache = NewCache("./", "instrumentsCache", "json")

// this map represents all of kaiko's available instruments per exchange
// we're keeping it alive as static variable as it is not heavy on ram
var instrumentsByExchange map[string][]string
var lastRefreshSec int64 = 0

// reloads the on heap map instrumentsByExchange (core data structure of the service)
func reloadInstrumentsMap(instruments []Instrument) {

	// here we translate the api output format (instruments with all their attributes)
	// into a frendlier format with what we are trying to achieve: a map of instruments slices by exchange
	// exampe of filling:
	// res := append(m["binance"], "BTCUSD")
	// m["bitfinex"] = []string{"BTC-USD"}
	// example of usage:
	// binanceInstruments = m["binance"]

	// first we instantiate/reset the map
	instrumentsByExchange = make(map[string][]string)

	// now we loop through the api results to match our map format
	for _, instrument := range instruments {
		if exchangeExists(instrument.ExchangeCode) {
			// key (exchange) already exists
			instrumentsByExchange[instrument.ExchangeCode] = append(instrumentsByExchange[instrument.ExchangeCode], instrument.ExchangePairCode)
		} else {
			// new key (exchange)
			instrumentsByExchange[instrument.ExchangeCode] = []string{instrument.ExchangePairCode}
		}
	}
	// once this is over, we can safely save our map in a persistent cache
	// so that we don't have to do this again if the server is restarted
	// or RAM is flushed
	if len(instrumentsByExchange) > 1 {
		// let's build the new cache from the fresh data
		err := instrumentsCache.writeSerialize(instrumentsByExchange)
		if err != nil {
			println("instruments cache could not be updated:", err)
		}
	}
}

// TODO: explore export possibilities for use in foreign packages
// returns slice of all instruments for a given exchange
func getInstrumentsForExchange(exchange string) ([]string, error) {

	nowUnix := time.Now().Unix()
	sinceRefres := nowUnix - lastRefreshSec
	if sinceRefres > CACHE_TIMEOUT_SECONDS {
		// check for refresh
		err := refreshInstrumentsByExchange()
		if err != nil {
			return nil, err
		}
		lastRefreshSec = nowUnix
	}

	// check if the exchange is available in our instrumentsByExchange map
	if val, ok := instrumentsByExchange[exchange]; ok {
		return val, nil
	} else {
		return nil, errors.New("no echange with such a name available")
	}
}

// TODO: explore export possibilities for use in foreign packages
// checks if the given exchange exists
func getExchanges() ([]string, error) {

	// check for refresh
	err := refreshInstrumentsByExchange()
	if err != nil {
		return nil, err
	}
	exchanges := make([]string, 0, len(instrumentsByExchange))
	for k := range instrumentsByExchange {
		exchanges = append(exchanges, k)
	}
	return exchanges, nil
}

// checks if the given exchange exists
func exchangeExists(exchange string) bool {
	_, ok := instrumentsByExchange[exchange]
	return ok
}

// checks if the given isntrument exists for the given exchange
func instrumentExistsStatus(exchange string, instrument string) int32 {

	// check if the exchange exists and if we get instruments for it
	instruments, err := getInstrumentsForExchange(exchange)
	if err != nil {
		// UNKNOWN exchange
		fmt.Printf("<< [0] unkown exchange\n")
		return 0
	}

	// now iterate over the exchange instruments to match with the instrument
	for _, ref := range instruments {
		if ref == instrument {
			// YES is exists
			fmt.Printf("<< [1] instrument available\n")
			return 1
		}
	}
	// NO it does not
	fmt.Printf("<< [2] instrument non available\n")
	return 2
}

// this function routes the retrieval of instruments either locally or through the api
func refreshInstrumentsByExchange() error {

	// check whether the cache is still usable
	if instrumentsCache.isFromLastSeconds(CACHE_TIMEOUT_SECONDS) {
		// if it is, but our map is empty, let's load from it
		if len(instrumentsByExchange) < 1 {
			err := fetchInstrumentsFromCache()
			if err != nil {
				return err
			}
		}
	} else {
		// the cache was outdated, let's try to hit the api
		fromRest, err := fetchInstrumentsRest()
		if err != nil {
			return err
		}
		// if we successfully retrieved at least 1 instrument
		if len(fromRest) > 1 {
			// build the map from this raw api data
			// and repopulate the cache
			reloadInstrumentsMap(fromRest)
			if err != nil {
				println("instruments cache could not be updated:", err)
			}
		}
	}
	return nil
}

// this function tries to retrieve the cached instruments data
func fetchInstrumentsFromCache() error {

	// let's check if the cache is fresh enough for us to load it straight
	if instrumentsCache.isFromLastSeconds(CACHE_TIMEOUT_SECONDS) {
		// fresh enough, load it
		instrumentsCache.readDeserialize(&instrumentsByExchange)
		return nil
	}
	// needs refresh
	return errors.New("instruments cache outdated")
}

func fetchInstrumentsRest() ([]Instrument, error) {

	// let's buit a GET request to match the endpoint spec
	req, err := http.NewRequest("GET", INSTRUMENTS_URL, nil)
	if err != nil {
		return nil, err
	}

	// let's instantiate a Client to enventually helps us with:
	// control over HTTP client headers
	// redirect policy ...
	client := &http.Client{}

	// send the request via a client, retrieve response
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// close res.Body when out of scope
	defer res.Body.Close()

	// fill the formattedRes with the data from the JSON
	var formattedRes InstrumentsResponse

	// read though the stream of JSON data and map in to our InstrumentsResponse entity
	if err := json.NewDecoder(res.Body).Decode(&formattedRes); err != nil {
		return nil, err
	}
	// check if the parsed data has at least 1 pair
	if len(formattedRes.Data) > 0 {
		return formattedRes.Data, nil
	}
	return nil, errors.New("No instruments returned")
}
