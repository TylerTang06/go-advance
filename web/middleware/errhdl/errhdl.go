package errhdl

import "github.com/TylerTang06/go-advance/web"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: make(map[int][]byte, 64),
	}
}

func (m *MiddlewareBuilder) RegisterError(code int, resp []byte) *MiddlewareBuilder {
	m.resp[code] = resp
	return m
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = resp
			}
		}
	}
}
