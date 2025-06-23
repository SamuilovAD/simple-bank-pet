package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
)

const (
	grpcUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader = "x-forwarded-for"
	userAgentHeader     = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	metadataStruct := &Metadata{}
	if requestMetadataMap, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("requestMetadataMap: %v\n", requestMetadataMap)
		if userAgents := requestMetadataMap.Get(grpcUserAgentHeader); len(userAgents) > 0 {
			metadataStruct.UserAgent = userAgents[0]
		}
		if userAgents := requestMetadataMap.Get(userAgentHeader); len(userAgents) > 0 {
			metadataStruct.UserAgent = userAgents[0]
		}
		if clientIps := requestMetadataMap.Get(xForwardedForHeader); len(clientIps) > 0 {
			metadataStruct.ClientIP = clientIps[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		metadataStruct.ClientIP = p.Addr.String()
	}

	return metadataStruct
}
