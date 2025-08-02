package reqctx

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
