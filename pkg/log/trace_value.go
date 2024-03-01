package log

import (
	"context"
)

const (
	requestId = "X-Request-Id"
	traceId   = "X-B3-Traceid"
	spanId    = "X-B3-Spanid"
)

func RequestID() Valuer {
	return func(ctx context.Context) interface{} {
		if ctx == nil {
			return ""
		}
		return ctx.Value(requestId)
	}
}

func RequestIDWithName(name string) Valuer {
	return func(ctx context.Context) interface{} {
		if ctx == nil {
			return ""
		}
		return ctx.Value(name)
	}
}

func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		if ctx == nil {
			return ""
		}
		return ctx.Value(traceId)
	}
}

func TraceIDWithName(name string) Valuer {
	return func(ctx context.Context) interface{} {
		if ctx == nil {
			return ""
		}
		return ctx.Value(name)
	}
}

func SpanID() Valuer {
	return func(ctx context.Context) interface{} {
		if ctx == nil {
			return ""
		}
		return ctx.Value(spanId)
	}
}

func SpanIDWithName(name string) Valuer {
	return func(ctx context.Context) interface{} {
		if ctx == nil {
			return ""
		}
		return ctx.Value(name)
	}
}
