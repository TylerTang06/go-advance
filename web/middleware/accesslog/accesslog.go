package accesslog

import (
	"encoding/json"
	"github.com/TylerTang06/go-advance/web"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(accessLog string)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

type accessLog struct {
	Host       string
	Route      string
	HttpMethod string `json:"http_method"`
	Path       string
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					Path:       ctx.Req.URL.Path,
					HttpMethod: ctx.Req.Method,
				}
				val, _ := json.Marshal(l)
				m.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}
