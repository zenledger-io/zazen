package log

import "context"

const (
	contextUUIDKey contextKey = "uuid"
)

type contextKey string

func ContextUUID(ctx context.Context) string {
	uuid, _ := ctx.Value(contextUUIDKey).(string)
	return uuid
}

func ContextWithUUID(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, contextUUIDKey, uuid)
}
