package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Print(md)

		userAgents := md.Get(grpcGatewayUserAgentHeader)
		if len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		userAgents = md.Get(userAgentHeader)
		if len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		clientIPs := md.Get(xForwardedForHeader)
		if len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	peer, ok := peer.FromContext(ctx)
	if ok {
		mtdt.ClientIP = peer.Addr.String()
	}

	return mtdt
}
