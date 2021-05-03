package middleware

import (
	"kisgateway/kishttp"
	"kisgateway/serverlib/logx"
	"net/http/httputil"
	"net/url"
)

func ReverseProxyMiddleWare() kishttp.HandlerFunc {
	return func(c *kishttp.Context) {
		logx.Info("http://127.0.0.1:2003/base")
		backURL, err := url.Parse("http://127.0.0.1:2003/base")
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(backURL)

		proxy.ServeHTTP(c.Writer, c.Req)
	}
}