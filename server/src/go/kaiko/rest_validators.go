package kaiko

import "time"

// these were generated with https://mholt.github.io/json-to-go/
// used to parse the http request on https://reference-data-api.kaiko.io/v1/instruments
type Instrument struct {
	KaikoLegacyExchangeSlug string    `json:"kaiko_legacy_exchange_slug"`
	TradeStartTime          time.Time `json:"trade_start_time"`
	TradeEndTime            time.Time `json:"trade_end_time"`
	Code                    string    `json:"code"`
	ExchangeCode            string    `json:"exchange_code"`
	ExchangePairCode        string    `json:"exchange_pair_code"`
	BaseAsset               string    `json:"base_asset"`
	QuoteAsset              string    `json:"quote_asset"`
	KaikoLegacySymbol       string    `json:"kaiko_legacy_symbol"`
	Class                   string    `json:"class"`
	TradeStartTimestamp     int64     `json:"trade_start_timestamp"`
	TradeEndTimestamp       int64     `json:"trade_end_timestamp"`
	TradeCount              int       `json:"trade_count"`
	TradeCompressedSize     int       `json:"trade_compressed_size"`
}

type InstrumentsResponse struct {
	Result string       `json:"result"`
	Data   []Instrument `json:"data"`
}
