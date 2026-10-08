package interceptor

import (
	"context"
	"strings"

	"github.com/MaksimCpp/pkg/authcontext"

	authpb "github.com/MaksimCpp/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	RoleAdmin   string = "admin"
	// RoleUser    string = "user"
	// RoleCourier string = "courier"
)

func AuthInterceptor(
	authClient authpb.AuthServiceClient,
	publicMethods map[string]struct{},
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {

		if _, ok := publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

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

		tokenRequest := authpb.ValidateRequest{
			AccessToken: token,
		}

		tokenResponse, err := authClient.ValidateAccessToken(ctx, &tokenRequest)
		if err != nil {
			return nil, status.Error(
				codes.Unauthenticated, "Invalid access token",
			)
		}

		ctx = context.WithValue(
			ctx,
			authcontext.UserIDKey,
			tokenResponse.UserId,
		)

		ctx = context.WithValue(
			ctx,
			authcontext.RoleKey,
			tokenResponse.Role,
		)
		return handler(ctx, req)
	}
}

func RoleInterceptor(methodsForAdmin map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		if _, ok := methodsForAdmin[info.FullMethod]; !ok {
			return handler(ctx, req)
		}

		role, ok := authcontext.GetRole(ctx)

		if !ok {
			return nil, status.Error(
				codes.Unauthenticated, "Unauthenticated",
			)
		}

		if role != RoleAdmin {
			return nil, status.Error(
				codes.PermissionDenied, "Permission denied",
			)
		}

		return handler(ctx, req)
	}
}
