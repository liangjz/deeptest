package websocketHelper

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	_i118Utils "github.com/aaronchen2k/deeptest/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"strings"
)

var (
	wsConn *neffos.Conn
)

func SendExecMsg(msg string, log domain.ExecLog, wsMsg *websocket.Message) {
	msg = strings.TrimSpace(msg)
	resp := _domain.WsResp{Msg: msg, Data: log}

	bytes, _ := json.Marshal(resp)
	mqData := _domain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}

	msg = strings.ReplaceAll(strings.TrimSpace(msg), `%`, `%%`)
	if wsMsg != nil {
		logUtils.Infof(_i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room, msg))
		PubMsg(mqData)
	} else {
		logUtils.Infof(msg)
	}
}

func SendExecResult(category consts.WsMsgCategory, data interface{}, wsMsg *websocket.Message) {
	resp := _domain.WsResp{Category: category, Data: data}
	if data != nil {
		resp.Data = data
	}
	bytes, _ := json.Marshal(resp)
	mqData := _domain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}

	if wsMsg != nil {
		logUtils.Infof(_i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room, category))
		PubMsg(mqData)
	} else {
		logUtils.Infof(string(bytes))
	}
}

func Broadcast(namespace, room, event string, content string) {
	wsConn.Server().Broadcast(nil, websocket.Message{
		Namespace: namespace,
		Room:      room,
		Event:     event,
		Body:      []byte(content),
	})
}

func SetConn(conn *neffos.Conn) {
	wsConn = conn
}

type PrefixedLogger struct {
	Prefix string
}

func (s *PrefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.Prefix, msg)
}
