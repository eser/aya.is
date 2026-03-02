package mcp

import (
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "aya-mcp"
	serverVersion = "1.0.0"
)

func RegisterMCPRoutes(
	routes *httpfx.Router,
	profileService *profiles.Service,
	storyService *stories.Service,
) {
	server := mcp.NewServer(
		&mcp.Implementation{ //nolint:exhaustruct // external SDK type
			Name:    serverName,
			Version: serverVersion,
		},
		nil,
	)

	registerProfileTools(server, profileService)
	registerStoryTools(server, storyService)
	registerNewsTools(server, storyService)

	handler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{ //nolint:exhaustruct // external SDK type
			Stateless: true,
		},
	)

	// Register with explicit HTTP methods to avoid conflict with OPTIONS wildcard
	routes.GetMux().Handle("GET /mcp", handler)
	routes.GetMux().Handle("POST /mcp", handler)
}
