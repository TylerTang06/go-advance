package cookie

import (
	"github.com/TylerTang06/go-advance/web/session"
	"net/http"
)

var _ session.Propagator = &Propagator{}

type PropagatorOption func(propagator *Propagator)

func WithCookieOption(opt func(c *http.Cookie)) PropagatorOption {
	return func(propagator *Propagator) {
		propagator.cookieOpt = opt
	}
}

type Propagator struct {
	cookieName string
	cookieOpt  func(cookie *http.Cookie)
}

func NewPropagator(cookieName string, cookieOpts ...PropagatorOption) *Propagator {
	return &Propagator{
		cookieName: cookieName,
		cookieOpt:  func(cookie *http.Cookie) {},
	}
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
	}
	p.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	}
	p.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}
