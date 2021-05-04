package middleware

import (
	"kisgateway/kishttp"
	"kisgateway/serverlib/jwt"
	"strings"
)

func JwtAuthMiddleware() kishttp.HandlerFunc {
	return func(c *kishttp.Context) {
		service := getService(c)
		if service.AccessControl.OpenAuth != 1 {
			return
		}
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		//fmt.Println("token",token)
		if token != "" {
			claims, err := jwt.JwtDecode(token)
			if err != nil {
				c.String( 2003, "token err")
				c.Abort()
				return
			}
			//fmt.Println("claims.Issuer",claims.Issuer)
			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					c.Next()
					return
				}
			}
		}
		c.String( 2003, "not match valid app")
		c.Abort()
		return
	}
}
