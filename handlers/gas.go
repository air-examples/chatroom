package handlers

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
			c, ok := req.Cookies["name"]
			if !ok {
				air.ERROR("name not found in cookie")
				return Error(
					res,
					401,
					"name not found in cookie",
				)
			}
			c.Expires = time.Now().Add(7 * 24 * time.Hour)
			user, ok := models.GetUser(c.Value)
			if !ok {
				air.ERROR("name not found in cache")
				return Error(
					res,
					401,
					"name not found in cache",
				)
			}
			req.Values["name"] = user.Name
			res.Cookies["name"] = c
			return next(req, res)
		}
	}
}
