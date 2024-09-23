package service

import (
	"fmt"
	agentExec "github.com/aaronchen2k/deeptest/internal/agent/exec"
	execUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	stringUtils "github.com/aaronchen2k/deeptest/pkg/lib/string"
	"github.com/kataras/iris/v12/websocket"
	"github.com/savsgio/gotils/strings"
	"runtime/debug"
)

func StartExec(req agentExec.WsReq, wsMsg *websocket.Message) (err error) {
	act := req.Act
	execUuid := getExecUuid(req)
	if execUuid == "" {
		logUtils.Info("****** execUuid is empty")
		logUtils.Infof("%v", req)
		return
	}

	// stop exec
	if act == consts.ExecStop {
		StopExec(execUuid, wsMsg)
		return
	}

	// is running
	ctx, _ := agentExec.GetExecCtx(execUuid)
	if ctx != nil && (strings.Include([]string{consts.ExecScenario.String(), consts.ExecPlan.String(), consts.ExecCase.String()}, act.String())) {
		execUtils.SendAlreadyRunningMsg(wsMsg)
		return
	}

	// exec task
	go func() {
		defer errDefer(wsMsg)

		agentExec.InitUserExecContext(execUuid)

		if act == consts.ExecInterface {
			req.InterfaceExecReq.TenantId = req.TenantId
			RunInterface(&req.InterfaceExecReq, req.LocalVarsCache, wsMsg)

		} else if act == consts.ExecScenario {
			req.ScenarioExecReq.TenantId = req.TenantId
			RunScenario(&req.ScenarioExecReq, req.LocalVarsCache, wsMsg)

		} else if act == consts.ExecPlan {
			req.PlanExecReq.TenantId = req.TenantId
			RunPlan(&req.PlanExecReq, req.LocalVarsCache, wsMsg)

		} else if act == consts.ExecCase {
			req.CasesExecReq.TenantId = req.TenantId
			RunCases(&req.CasesExecReq, req.LocalVarsCache, wsMsg)

		} else if act == consts.ExecMessage {
			req.MessageReq.TenantId = req.TenantId
			RunMessage(&req.MessageReq, wsMsg)

		} else if stringUtils.FindInArr(act.String(), consts.WebscoketAtions) {
			req.WebsocketExecReq.TenantId = req.TenantId
			RunWebsocket(act, &req.WebsocketExecReq, req.LocalVarsCache, wsMsg)

		}

		if !stringUtils.FindInArr(act.String(), consts.WebscoketAtions) { // keep context for websocket testing
			agentExec.CancelExecCtx(execUuid)
		}
	}()

	return
}

func getExecUuid(req agentExec.WsReq) (ret string) {
	if req.InterfaceExecReq.ExecUuid != "" {
		ret = req.InterfaceExecReq.ExecUuid

	} else if req.ScenarioExecReq.ExecUuid != "" {
		ret = req.ScenarioExecReq.ExecUuid

	} else if req.PlanExecReq.ExecUuid != "" {
		ret = req.PlanExecReq.ExecUuid

	} else if req.CasesExecReq.ExecUuid != "" {
		ret = req.CasesExecReq.ExecUuid

	} else if req.WebsocketExecReq.Room != "" {
		ret = req.WebsocketExecReq.Room

	}

	return
}

func StopExec(execUuid string, wsMsg *websocket.Message) (err error) {
	agentExec.CancelExecCtx(execUuid)
	execUtils.SendCancelMsg(wsMsg)

	return
}

func errDefer(wsMsg *websocket.Message) {
	err := recover()

	if err != nil {
		s := string(debug.Stack())
		fmt.Printf("err=%v, stack=%s\n", err, s)

		execUtils.SendErrorMsg(err, consts.Processor, wsMsg)
	}
}