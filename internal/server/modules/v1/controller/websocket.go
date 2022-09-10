package controller

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	execHelper "github.com/aaronchen2k/deeptest/internal/server/modules/v1/helper/exec"
	websocketHelper "github.com/aaronchen2k/deeptest/internal/server/modules/v1/helper/websocket"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	_i118Utils "github.com/aaronchen2k/deeptest/pkg/lib/i118"
	_logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/kataras/iris/v12/websocket"
)

const (
	result = "result"
	outPut = "output"
)

var (
	ch chan int
)

type WebSocketCtrl struct {
	Namespace         string
	*websocket.NSConn `stateless:"true"`

	ScenarioExecService *service.ExecScenarioService `inject:""`
}

func NewWsCtrl() *WebSocketCtrl {
	inst := &WebSocketCtrl{Namespace: serverConsts.WsDefaultNameSpace}
	return inst
}

func (c *WebSocketCtrl) OnNamespaceConnected(wsMsg websocket.Message) error {
	websocketHelper.SetConn(c.Conn)

	_logUtils.Infof(_i118Utils.Sprintf("ws_namespace_connected", c.Conn.ID(), wsMsg.Room))

	resp := _domain.WsResp{Msg: "from server: connected to websocket"}
	bytes, _ := json.Marshal(resp)
	mqData := _domain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	websocketHelper.PubMsg(mqData)

	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design.
func (c *WebSocketCtrl) OnNamespaceDisconnect(wsMsg websocket.Message) error {
	_logUtils.Infof(_i118Utils.Sprintf("ws_namespace_disconnected", c.Conn.ID()))

	resp := _domain.WsResp{Msg: fmt.Sprintf("ws_connected")}
	bytes, _ := json.Marshal(resp)
	mqData := _domain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	websocketHelper.PubMsg(mqData)

	return nil
}

// OnChat This will call the "OnVisit" event on all clients, including the current one, with the 'newCount' variable.
func (c *WebSocketCtrl) OnChat(wsMsg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)

	_logUtils.Infof("WebSocket OnChat: remote address=%s, room=%s, msg=%s", ctx.RemoteAddr(), wsMsg.Room, string(wsMsg.Body))

	req := domain.WsReq{}
	err = json.Unmarshal(wsMsg.Body, &req)
	if err != nil {
		execHelper.SendErrorMsg(req.Id, wsMsg)
		return
	}

	act := req.Act

	if act == consts.ExecStop {
		if ch != nil {
			if !execHelper.GetRunning() {
				ch = nil
			} else {
				ch <- 1
				ch = nil
			}
		}

		c.ScenarioExecService.CancelAndSendMsg(req.Id, wsMsg)

		return
	}

	if execHelper.GetRunning() && (act == consts.ExecStart) { // already running
		execHelper.SendAlreadyRunningMsg(req.Id, wsMsg)
		return
	}

	ch = make(chan int, 1)
	go func() {
		c.ScenarioExecService.ExecScenario(req.Id, &wsMsg)

		//c.Complete(wsMsg)
	}()

	return
}
