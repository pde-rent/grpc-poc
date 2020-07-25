package kaiko

import (
	fmt "fmt"
	"testing"
)

// main test function
func TestGrpcExistsExchange(t *testing.T) {

	// test watchers
	var existingRunCount int = 0
	var fakeRunCount int = 0
	var runCount int = 0
	var errCount int = 0
	var averageLatency int
	var maxLatency int = 0
	var minLatency int = 1000000000

	// run int tests on existing instruments proc calls
	println("testing...")

	var res bool
	var success string
	var ms int
	var err error
	var exchange string
	var instrument string

	for i := 1; i < (INTEGRATION_TEST_RUNS_EXISTING + INTEGRATION_TEST_RUNS_FAKE + 1); i++ {

		if runCount < INTEGRATION_TEST_RUNS_FAKE {
			// testing existing instruments / echanges combinations
			exchange, instrument = getExistingInstrument()
			res, ms, err = testExistsSingle(exchange, instrument)
			existingRunCount++
			// the RPC should return true
			if res == true {
				success = "✔️"
			} else {
				success = "❌"
			}
		} else {
			// testing existing instruments / echanges combinations
			exchange, instrument = getFakeInstrument()
			res, ms, err = testExistsSingle(exchange, instrument)
			fakeRunCount++
			// the RPC should return false
			if res == false {
				success = "✔️"
			} else {
				success = "❌"
			}
		}
		// update global counts
		runCount++
		if err != nil {
			errCount++
		}
		// print the test log
		fmt.Printf("%d] \tout RPC >> {%s|%s} \t<< %s\n", runCount, exchange, instrument, success)
		// update latency watchers
		averageLatency = (averageLatency*(runCount-1) + ms) / runCount
		if ms < minLatency {
			minLatency = ms
		}
		if ms > maxLatency {
			maxLatency = ms
		}
	}
	// print integration test result
	fmt.Printf(`
		runCount: %d
		existingRunCount: %d
		fakeRunCount: %d
		errCount: %d
		averageLatency: %dms
		maxLatency: %dms
		minLatency: %dms
	`, runCount, existingRunCount, fakeRunCount, errCount, averageLatency, maxLatency, minLatency)
}
