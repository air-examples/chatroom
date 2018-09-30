package handlers

import (
	"errors"
	"strings"

	"github.com/air-examples/chatroom/gas"
	"github.com/air-examples/chatroom/models"

	"github.com/aofei/air"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	users = cmap.New()
)

func init() {
	air.GET("/chat", chatPageHandler, gas.Auth)
	air.POST("/api/name", chatNameHandler)
	air.GET("/socket", socketHandler, gas.Auth)
}

func chatPageHandler(req *air.Request, res *air.Response) error {
	return res.Render(req.Values, "chat.html", "base.html")
}

func chatNameHandler(req *air.Request, res *air.Response) error {
	air.INFO("in chat name handler")
	name := strings.TrimSpace(req.Params["name"])
	k := models.GetAuthKey(name)
	if _, ok := models.GetUser(k); ok {
		air.ERROR("duplicate name", map[string]interface{}{
			"name":    name,
			"request": req,
		})
		return Error(400, errors.New("duplicate name"))
	}
	models.SetUser(models.NewUser(name))
	res.Cookies["name"] = &air.Cookie{
		Name:  "name",
		Value: k,
		Path:  "/",
	}
	return Success(res, "")
}

func socketHandler(req *air.Request, res *air.Response) error {
	c, err := res.UpgradeToWebSocket()
	if err != nil {
		air.ERROR("upgrade to websocket error", map[string]interface{}{
			"request": req,
			"error":   err.Error(),
		})
		return Error(500, err)
	}
	defer c.Close()

	name := req.Params["name"]
	if _, ok := users.Get(name); ok {
		air.ERROR("duplicate name", map[string]interface{}{
			"req":  req,
			"name": name,
		})
		return Error(400, errors.New("duplicate name"))
	}

	me := newSocketManager(name)
	users.Set(name, me)

	go func() {
		for {
			if t, b, err := c.ReadMessage(); err == nil {
				switch t {
				case air.WebSocketMessageTypeText:
					me.SendMsg(newMsg(name, t, b))
				case air.WebSocketMessageTypeBinary:
				case air.WebSocketMessageTypeConnectionClose:
					me.Close()
					return
				}
			} else {
				air.ERROR("socket msg error",
					map[string]interface{}{
						"type":    t,
						"err":     err.Error(),
						"content": string(b),
					})
				me.Close()
				return
			}
		}
	}()

	for {
		select {
		case <-me.newMsg:
			if v, err := me.msg.Marshal(); err == nil {
				err = c.WriteMessage(me.msg.MType, v)
				if err != nil {
					air.ERROR("send socket msg error",
						map[string]interface{}{
							"content": me.msg,
							"to":      me.name,
							"err":     err,
						})
					me.Close()
				}
				me.mu.Unlock()
			}
		case <-me.shutdown:
			break
		}
	}

	return nil
}
