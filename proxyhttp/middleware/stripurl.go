package middleware

import (
	"kisgateway/kishttp"
	"kisgateway/serverlib/logx"
)

func StripUrlMiddleware() kishttp.HandlerFunc {
	return func(c *kishttp.Context) {
		logx.Info("StripUrlMiddleware")
		c.Req.Header.Add("stripUrl", "test")
		c.Writer.Header().Set("stripUrl", "test")
	}
}
