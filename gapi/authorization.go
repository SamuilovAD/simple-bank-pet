package gapi

import (
	"fmt"
	"github.com/SamuilovAD/simple-bank-pet/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationBearer    = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata not found")

	}
	values := md.Get(authorizationHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("authorization header not found")
	}

	authHeader := values[0]
	fields := strings.Split(authHeader, " ")
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}
	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token %s", err)
	}
	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
