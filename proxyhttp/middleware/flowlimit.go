package middleware

import (
	"kisgateway/kishttp"
	"kisgateway/serverlib/limit"
	"net/http"
)

func FlowLimitMiddleware() kishttp.HandlerFunc {
	return func(c *kishttp.Context) {
		service := getService(c)
		serviceControl := service.AccessControl
		//服务端限流
		if serviceControl.ServiceFlowLimit != 0 {
			checkFlowLimit(c, service.ServerInfo.ServiceName, float64(serviceControl.ServiceFlowLimit))
		}
		//客户端限流
		if serviceControl.ClientipFlowLimit != 0 {
			checkFlowLimit(c, service.ServerInfo.ServiceName+c.ClientIP(), float64(serviceControl.ClientipFlowLimit))
		}
	}
}

func checkFlowLimit(c *kishttp.Context, serviceKey string, qps float64) {
	l, err := limit.FlowLimiterHandler.GetLimiter(serviceKey, qps)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		c.Abort()
		return
	}
	if !l.Allow() {
		c.String(http.StatusBadRequest, "")
		c.Abort()
		return
	}
}
