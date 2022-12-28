package gapi

import (
	"context"
	"log"

	// _ "google.golang.org/grpc/internal/metadata"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	httpGatewayUserAgentHeader = "grpcgateway-user-agent"
	grpcGatewayUserAgentHeader = "user-agent"

	httpGatewayUserClientIPHeader = "x-forwarded-for"
	grpcGatewayUserClientIPHeader = ""
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

//gRPC client IP address is not stored in the metadata it seems. Though for http requests and grpc requests, the user agent can be found in the metadata.
//for http requests the client IP address can be found.

//to find the client IP address in grpc requests, we have to find it somewhere in the context, view code below

//Keep in mind, only a grpc gateway can accept both http and grpc requests
func (server *Server) extractMetadata(ctx context.Context) *Metadata {

	finalMetadata := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md: %+v\n", md)

		if userAgents := md.Get(httpGatewayUserAgentHeader); len(userAgents) > 0 {
			finalMetadata.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			finalMetadata.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(httpGatewayUserClientIPHeader); len(clientIPs) > 0 {
			finalMetadata.ClientIP = clientIPs[0]
		}
	}

	//code to find client IP address from the the context in a grpc request
	if p, ok := peer.FromContext(ctx); ok {
		finalMetadata.ClientIP = p.Addr.String()
	}

	return finalMetadata
}
