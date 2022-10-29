package handler

import (
	"encoding/json"
	"fmt"
	execDomain "github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/agent/service"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/websocket"
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

	ExecService *service.ExecService `inject:""`
}

func NewWsCtrl() *WebSocketCtrl {
	inst := &WebSocketCtrl{Namespace: consts.WsDefaultNameSpace}
	return inst
}

func (c *WebSocketCtrl) OnNamespaceConnected(wsMsg websocket.Message) error {
	websocketHelper.SetConn(c.Conn)

	_logUtils.Infof(_i118Utils.Sprintf("ws_namespace_connected", c.Conn.ID(), wsMsg.Room))

	resp := _domain.WsResp{Msg: "from agent: connected to websocket"}
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

	resp := _domain.WsResp{Msg: "from agent: disconnected to websocket"}
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
		sendErr(err, &wsMsg)
		return
	}

	act := req.Act

	if act == consts.ExecStop {
		if ch != nil {
			if !exec.GetRunning() {
				ch = nil
			} else {
				ch <- 1
				ch = nil
			}
		}

		c.ExecService.CancelAndSendMsg(req.Id, wsMsg)

		return
	}

	if exec.GetRunning() && (act == consts.ExecStart) { // already running
		exec.SendAlreadyRunningMsg(req.Id, wsMsg)
		return
	}

	if act == consts.ExecScenario {
		ch = make(chan int, 1)
		go func() {
			c.ExecService.ExecScenario(&req.ExecReq, &wsMsg)
		}()
	}

	return
}

func sendErr(err error, wsMsg *websocket.Message) {
	root := execDomain.Result{
		ID:      -1,
		Name:    "执行失败",
		Summary: fmt.Sprintf("错误：%s", err.Error()),
	}
	exec.SendExecMsg(root, wsMsg)

	result := execDomain.Result{
		ID:       -2,
		ParentId: -1,
		Name:     "执行失败",
		Summary:  fmt.Sprintf("错误：%s", err.Error()),
	}
	exec.SendExecMsg(result, wsMsg)
}