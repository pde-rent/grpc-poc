package kaiko

import (
	// "log"
	"context"
	fmt "fmt"
)

// TODO: replace the fmt.Printf and printl with proper logs in logfiles

// the Server represents the gRPC server
// it implements protobuf services (endpoint as handlers)
type Server struct {
	UnimplementedKaikoServer
}

var existsCallCount = 0

// ExistsHandler generates response to a Ping request
func (s *Server) Exists(ctx context.Context, req *ExistsRequest) (*ExistsResponse, error) {
	existsCallCount++
	fmt.Printf("%d] \tin RPC >> {%s|%s} \t", existsCallCount, req.ExchangeCode, req.ExchangePairCode)
	// check with the updated Kaiko available echangePairs list for a match
	// the service function is implemented in the rest_consumer.go file
	// TODO: move this to a standalone service file for cleaner structure
	existCode := instrumentExistsStatus(req.ExchangeCode, req.ExchangePairCode)
	return &ExistsResponse{Exists: ExistsResponse_Exists(existCode)}, nil
}

// TODO: following the model above we can here implement as many RPC calls
// as required by the application
// a good practice would be to disociate the services and the front facing handlers
