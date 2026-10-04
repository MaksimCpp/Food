package authcontext

import (
	"context"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const RoleKey contextKey = "role"

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)

	if !ok {
		return -1, false
	}

	return userID, true
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)

	if !ok {
		return "", false
	}

	return role, true
}