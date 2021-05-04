package middleware

import (
	"kisgateway/kishttp"
	"kisgateway/serverlib/store"
	"net/http"
	"strings"
)

func FindServiceMiddle() kishttp.HandlerFunc {
	return func(c *kishttp.Context) {
		host := c.Req.Host
		path := c.Req.URL.Path
		servicesList := store.Store.HttpServices
		for _, service := range servicesList {
			//域名类型
			if service.ServerInfo.LoadType == 0 && service.HttpRule.Rule == host {
				c.Set("service", service)
				return
			}
			//路径类型
			if service.ServerInfo.LoadType == 1 &&
				strings.HasPrefix(service.HttpRule.Rule, path) {
				c.Set("service", service)
				return
			}
		}
		c.String(http.StatusBadRequest, "bad request")
		c.Abort()
	}
}
