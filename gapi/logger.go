package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This gRPC unary interceptor just incoporates the function signature of the UnaryInterceptor declared
// in the gRPC source code
func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	//To get the duration of the request. We record the time the request started
	startTime := time.Now()

	result, err := handler(ctx, req)

	//Then we get the time since the request started
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if status, ok := status.FromError(err); ok {
		statusCode = status.Code()
	}

	//This is done in a bid to differentiate between info and error events
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Str("status_text", statusCode.String()).
		Dur("request_duration", duration).
		Int("status_code", int(statusCode)).
		Msg("recieved a grpc request!")

	return result, err
}
