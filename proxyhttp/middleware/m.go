package middleware

import (
	"kisgateway/kishttp"
	"kisgateway/serverlib/store/gateway"
	"net/http"
)

func getService(c *kishttp.Context) gateway.ServiceInfo {
	serviceI, ok := c.Get("service")
	if !ok {
		c.String(http.StatusBadRequest, "")
		c.Abort()
		return gateway.ServiceInfo{}
	}
	return serviceI.(gateway.ServiceInfo)
}