package gas

import (
	"errors"
	"time"

	"github.com/air-examples/chatroom/common"
	"github.com/air-examples/chatroom/utils"
	"github.com/sheng/air"
)

func AuthHandler() air.Gas {
	return func(next air.Handler) air.Handler {
		return func(req *air.Request, res *air.Response) error {
			// if air.DebugMode {
			// 	air.INFO("debug mode, pass")
			// 	return nil
			// }
			name := ""
			cookie := &air.Cookie{}
			for _, c := range req.Cookies {
				if c.Name == common.AuthCookie {
					name = c.Value
					c.Expires = time.Now().Add(7 * 24 * time.Hour)
					cookie = c
					break
				}
			}
			if name == "" {
				air.ERROR("name not found in cookie")
				return utils.Error(401,
					errors.New("name not found"))
			}
			req.Params["name"] = name
			res.Cookies = append(res.Cookies, cookie)
			return next(req, res)
		}
	}
}
