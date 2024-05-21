package api

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestExtractMetadata(t *testing.T) {
	testCases := []struct {
		name     string
		headers  metadata.MD
		peerAddr net.IP
		expected Metadata
	}{
		{
			name: "ExtractMetadataFromGateway",
			headers: metadata.Pairs(
				grpcGatewayUserAgentHeader, "grpc-gateway-user-agent-value",
				xForwardedForHeader, "127.0.0.1",
			),
			expected: Metadata{
				UserAgent: "grpc-gateway-user-agent-value",
				ClientIP:  "127.0.0.1",
			},
		},
		{
			name: "ExtractMetadataFromGrpc",
			headers: metadata.Pairs(
				userAgentHeader, "user-agent-value",
			),
			peerAddr: net.ParseIP("127.0.0.1"),
			expected: Metadata{
				UserAgent: "user-agent-value",
				ClientIP:  "127.0.0.1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := &Server{}
			ctx := context.Background()

			// Set metadata in context
			if len(tc.headers) > 0 {
				ctx = metadata.NewIncomingContext(ctx, tc.headers)
			}

			// Set peer information in context
			if tc.peerAddr != nil {
				p := &peer.Peer{
					Addr: &net.IPAddr{
						IP: tc.peerAddr,
					},
				}
				ctx = peer.NewContext(ctx, p)
			}

			// Execute the extractMetadata method
			result := server.extractMetadata(ctx)

			// Verify the extracted metadata
			assert.Equal(t, tc.expected.UserAgent, result.UserAgent)
			assert.Equal(t, tc.expected.ClientIP, result.ClientIP)
		})
	}
}
