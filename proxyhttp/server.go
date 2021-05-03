package proxyhttp

import (
	"kisgateway/kishttp"
	"kisgateway/proxyhttp/middleware"
	"kisgateway/serverlib/conf"
	"kisgateway/serverlib/logx"
	"net/http"
	"time"
)

type HttpServer struct {
}

func New() *HttpServer {
	return &HttpServer{
	}
}

func (* HttpServer) Start()  {
	r := setRoute()
	s := http.Server{
		Addr:           conf.GetConf("http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(conf.GetConfInt("http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(conf.GetConfInt("http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(conf.GetConfInt("http.max_header_bytes")),
	}

	if err := s.ListenAndServe(); err != nil {
		logx.Info("proxy http run err: %s", err.Error())
	}
}

func setRoute() *kishttp.Engine {
	e := kishttp.New()
	e.GET("ping", func(c *kishttp.Context) {
		c.String(http.StatusOK, "pong")
	})

	e.Use(
		middleware.FindServiceMiddle(),
		middleware.StripUrlMiddleware(),
		middleware.ReverseProxyMiddleWare(),
	)
	return e
}
