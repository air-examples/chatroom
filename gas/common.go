package gas

import (
	"github.com/air-examples/chatroom/utils"
	"github.com/sheng/air"
)

func PreRenderHandler() air.Gas {
	return func(next air.Handler) air.Handler {
		return func(req *air.Request, res *air.Response) error {
			if _, ok := req.Values["title"]; !ok {
				req.Values["title"] =
					req.LocalizedString("title")
			}
			if _, ok := req.Values["username"]; !ok {
				req.Values["username"] =
					req.LocalizedString("username")
			}
			if _, ok := req.Values["confirm"]; !ok {
				req.Values["confirm"] =
					req.LocalizedString("confirm")
			}
			if _, ok := req.Values["confirm your name"]; !ok {
				req.Values["confirm your name"] =
					req.LocalizedString("confirm your name")
			}
			if _, ok := req.Values["send"]; !ok {
				req.Values["send"] =
					req.LocalizedString("send")
			}
			return next(req, res)
		}
	}
}

func PreLoggerHandler() air.Gas {
	return func(next air.Handler) air.Handler {
		return func(req *air.Request, res *air.Response) error {
			air.INFO("receive request", utils.M{
				"request": req,
			})
			return next(req, res)
		}
	}
}
