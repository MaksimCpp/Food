package authcontext

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"
)

func UserID(ctx context.Context) (int64, bool) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return -1, false
	}

	userIDValues := md.Get("user_id")
	userID, err := strconv.ParseInt(userIDValues[0], 10, 64)

	if err != nil {
		return -1, false
	}

	return userID, true
}