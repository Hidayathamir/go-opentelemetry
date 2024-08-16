package oteltrace

import (
	"context"

	"github.com/erajayatech/go-opentelemetry/pkg/caller"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func RecordSpan(ctx context.Context) (context.Context, trace.Span) {
	c, ok := ctx.(*gin.Context)
	if ok {
		return otel.Tracer("").Start(c.Request.Context(), caller.FuncName(caller.WithSkip(1))) //nolint:spancheck
	}
	return otel.Tracer("").Start(ctx, caller.FuncName(caller.WithSkip(1))) //nolint:spancheck
}
