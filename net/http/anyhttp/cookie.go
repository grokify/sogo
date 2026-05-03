package anyhttp

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

type Cookie struct {
	Name  string
	Value string
}

func (c *Cookie) ToNetHTTP() *http.Cookie {
	cookie := http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	return &cookie
}

func (c *Cookie) ToFastHTTP() *fasthttp.Cookie {
	var cookie fasthttp.Cookie
	cookie.SetKey(c.Name)
	cookie.SetValue(c.Value)
	return &cookie
}
