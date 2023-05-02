package service

import (
	"github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/kataras/iris/v12/websocket"
)

type ExecPlanService struct {
	RemoteService       *RemoteService       `inject:""`
	ExecScenarioService *ExecScenarioService `inject:""`
}

func (s *ExecPlanService) ExecPlan(req *agentExec.PlanExecReq, wsMsg *websocket.Message) (err error) {
	consts.ServerUrl = req.ServerUrl
	consts.ServerToken = req.Token

	planExecObj := s.RemoteService.GetPlanToExec(req)

	// start msg
	execUtils.SendStartMsg(wsMsg)

	// execution
	var results = agentDomain.PlanExecResult{
		ID: req.PlanId,
	}
	for _, scenario := range planExecObj.Scenarios {
		session, _ := s.ExecScenarioService.Exec(&scenario, wsMsg)
		results.Scenarios = append(results.Scenarios, session.RootProcessor.Result)
	}

	// submit result
	report, _ := s.RemoteService.SubmitPlanResult(results, req.PlanId, req.ServerUrl, req.Token)
	execUtils.SendResultMsg(report, wsMsg)
	s.sendSubmitResult(req.PlanId, wsMsg)

	// end msg
	execUtils.SendEndMsg(wsMsg)

	return
}

func (s *ExecPlanService) CancelAndSendMsg(planId int, wsMsg websocket.Message) (err error) {
	execUtils.SendCancelMsg(wsMsg)
	return
}

func (s *ExecPlanService) sendSubmitResult(planId uint, wsMsg *websocket.Message) (err error) {
	result := agentDomain.PlanExecResult{
		ID:   planId,
		Name: "提交执行结果成功",
		//Summary:  fmt.Sprintf("错误：%s", err.Error()),
	}
	execUtils.SendExecMsg(result, wsMsg)

	return
}