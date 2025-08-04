package web

import "context"

type contextKey string

const RequestIDKey contextKey = "requestID"

func GetRequestID(ctx context.Context) string {
	requestID := ""
	if id, ok := ctx.Value(RequestIDKey).(string); ok && id != "" {
		requestID = id
	}
	return requestID
}

const UserIDKey contextKey = "userID"

func GetUserID(ctx context.Context) string {
	userID := ""
	if id, ok := ctx.Value(UserIDKey).(string); ok && id != "" {
		userID = id
	}
	return userID
}
