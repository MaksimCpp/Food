package deliverygrpc

import (
	"context"
	"strings"

	"github.com/MaksimCpp/auth/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(tokenService domain.TokenService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, status.Error(
				codes.Unauthenticated, "Unauthenticated",
			)
		}

		values := md.Get("authorization")

		if len(values) == 0 {
			return nil, status.Error(
				codes.Unauthenticated, "Unauthenticated",
			)
		}

		authHeader := values[0]

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return nil, status.Error(
				codes.Unauthenticated, "Unauthenticated",
			)
		}

		token := parts[1]

		if token == "" {
			return nil, status.Error(
				codes.Unauthenticated, "Missing token",
			)
		}

		userID, err := tokenService.ValidateAccessToken(token)
		if err != nil {
			return nil, status.Error(
				codes.Unauthenticated, "Invalid access token",
			)
		}

		ctx = context.WithValue(
			ctx,
			"user_id",
			userID,
		)
		return handler(ctx, req)
	}
}