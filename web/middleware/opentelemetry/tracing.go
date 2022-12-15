package opentelemetry

import (
	"github.com/TylerTang06/go-advance/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const defaultInstrumentationName = "github.com/TylerTang06/go-advance/web/middle/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))
			reqCtx, span := m.Tracer.Start(reqCtx, "unknown", trace.WithAttributes())

			span.SetAttributes(attribute.String("http.mehod", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.host", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			// span.End 执行之后，就意味着 span 本身已经确定无疑了，将不能再变化了
			defer span.End()
			ctx.Req = ctx.Req.WithContext(reqCtx)
			next(ctx)

			// 使用命中的路由来作为 span 的名字
			if ctx.MatchedRoute != "" {
				span.SetName(ctx.MatchedRoute)
			}
			span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
		}
	}
}
