package gas

import (
	"time"

	"github.com/air-examples/chatroom/models"
	"github.com/aofei/air"
)

var (
	Auth = AuthHandler()
)

func AuthHandler() air.Gas {
	return func(next air.Handler) air.Handler {
		return func(req *air.Request, res *air.Response) error {
			// if air.DebugMode {
			// 	air.INFO("debug mode, pass")
			// 	return next(req, res)
			// }
			name := ""
			cookie := &air.Cookie{}
			c, ok := req.Cookies["name"]
			if !ok {
				air.ERROR("name not found in cookie")
				if req.Method == "GET" {
					req.URL.Path = "/"
					res.StatusCode = 302
					return res.Redirect(req.URL.String())
				}
				return &air.Error{
					Code:    401,
					Message: "name not found in cookie",
				}
			}
			name = c.Value
			c.Expires = time.Now().Add(7 * 24 * time.Hour)
			cookie = c
			v, ok := models.GetUser(name)
			if !ok {
				air.ERROR("name not found in cache")
				if req.Method == "GET" {
					req.URL.Path = "/"
					res.StatusCode = 302
					return res.Redirect(req.URL.String())
				}
				return &air.Error{
					Code:    400,
					Message: "name not found in cache",
				}
			}
			req.Params["name"] = v.Name
			res.Cookies["name"] = cookie
			return next(req, res)
		}
	}
}
