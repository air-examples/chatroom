package handlers

import (
	"errors"
	"strings"

	"github.com/air-examples/chatroom/common"
	"github.com/air-examples/chatroom/gas"
	"github.com/air-examples/chatroom/models"
	"github.com/air-examples/chatroom/utils"
	"github.com/sheng/air"
)

var (
	users map[string]*SocketManager
)

func init() {
	air.GET("/chat", chatPageHandler, gas.PreRender, gas.Auth)
	air.POST("/api/name", chatNameHandler)
	air.GET("/socket", socketHandler, gas.Auth)
}

func chatPageHandler(req *air.Request, res *air.Response) error {
	return res.Render(req.Values, "chat.html", "base.html")
}

func chatNameHandler(req *air.Request, res *air.Response) error {
	name := strings.TrimSpace(req.Params["name"])
	k := models.GetAuthKey(name)
	if _, ok := models.GetUser(k); ok {
		air.ERROR("duplicate name", utils.M{
			"name":    name,
			"request": req,
		})
		return utils.Error(400, errors.New("duplicate name"))
	}
	models.SetUser(models.NewUser(name))
	res.Cookies = append(res.Cookies, &air.Cookie{
		Name:  common.AuthCookie,
		Value: k,
		Path:  "/",
	})
	return utils.Success(res, "")
}

func socketHandler(req *air.Request, res *air.Response) error {
	c, err := res.UpgradeToWebSocket()
	if err != nil {
		air.ERROR("upgrade to websocket error", utils.M{
			"request": req,
			"error":   err.Error(),
		})
		return utils.Error(500, err)
	}
	defer c.Close()

	name := req.Params["name"]
	if _, ok := users[name]; ok {
		air.ERROR("duplicate name", utils.M{
			"req":  req,
			"name": name,
		})
		return utils.Error(400, errors.New("duplicate name"))
	}

	me := newSocketManager(name)
	if users == nil {
		users = make(map[string]*SocketManager)
	}
	users[name] = me

	go func() {
		for {
			if t, b, err := c.ReadMessage(); err == nil {
				switch t {
				case air.WebSocketMessageTypeText:
					me.SendMsg(newMsg(name, t, b))
				case air.WebSocketMessageTypeBinary:
				case air.WebSocketMessageTypeConnectionClose:
					delete(users, name)
					me.Close()
					return
				}
			} else {
				air.ERROR("socket msg error", utils.M{
					"type":    t,
					"err":     err.Error(),
					"content": string(b),
				})
				delete(users, name)
				me.Close()
				return
			}
		}
	}()

	for {
		select {
		case <-me.newMsg:
			if v, err := me.m.Marshal(); err == nil {
				err = c.WriteMessage(me.m.Mtype, v)
				if err != nil {
					air.ERROR("send socket msg error",
						utils.M{
							"content": me.m,
							"to":      me.name,
							"err":     err,
						})
				}
			}
		case <-me.shutdown:
			break
		}
	}

	return nil
}