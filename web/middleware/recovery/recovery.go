package recovery

import "github.com/TylerTang06/go-advance/web"

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *web.Context)
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = []byte(m.ErrMsg)
					m.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}
