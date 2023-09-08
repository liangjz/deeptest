package service

import (
	agentDomain "github.com/aaronchen2k/deeptest/cmd/agent/v1/domain"
	agentExec "github.com/aaronchen2k/deeptest/internal/agent/exec"
	execUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	_httpUtils "github.com/aaronchen2k/deeptest/pkg/lib/http"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
)

func RunInterface(call agentDomain.InterfaceCall) (resultReq domain.DebugData, resultResp domain.DebugResponse, err error) {
	req := GetInterfaceToExec(call)

	agentExec.CurrDebugInterfaceId = req.DebugData.DebugInterfaceId
	agentExec.CurrScenarioProcessorId = 0 // not in a scenario

	agentExec.ExecScene = req.ExecScene

	// init context
	agentExec.InitDebugExecContext()
	agentExec.InitJsRuntime()

	originalReqUri, _ := PreRequest(&req.DebugData)
	agentExec.SetReqValueToGoja(req.DebugData.BaseRequest)
	agentExec.ExecPreConditions(req)
	req.DebugData.BaseRequest = agentExec.CurrRequest // update to the value changed in goja

	resultResp, err = RequestInterface(&req.DebugData)

	agentExec.SetRespValueToGoja(resultResp)
	agentExec.ExecPostConditions(req, resultResp)
	agentExec.UpdateResponseData()
	PostRequest(originalReqUri, &req.DebugData)
	resultResp = agentExec.CurrResponse

	// submit result
	err = SubmitInterfaceResult(req, resultResp, call.ServerUrl, call.Token)

	resultReq = req.DebugData

	return
}

func PreRequest(req *domain.DebugData) (originalReqUri string, err error) {
	// replace variables
	agentExec.ReplaceVariables(&req.BaseRequest, consts.InterfaceDebug)

	// gen url
	originalReqUri = agentExec.ReplacePathParams(req.Url, req.PathParams)

	notUseBaseUrl := execUtils.IsUseBaseUrl(req.UsedBy, req.ProcessorInterfaceSrc)
	if notUseBaseUrl {
		req.BaseRequest.Url = originalReqUri
	} else {
		req.BaseRequest.Url = _httpUtils.CombineUrls(req.BaseUrl, originalReqUri)
	}
	req.BaseRequest.FullUrlToDisplay = req.BaseRequest.Url
	logUtils.Info("requested url: " + req.BaseRequest.Url)

	return
}

func PostRequest(originalReqUri string, req *domain.DebugData) (err error) {
	req.BaseRequest.Url = originalReqUri // rollback for saved to db

	return
}

func RequestInterface(req *domain.DebugData) (ret domain.DebugResponse, err error) {
	// send request
	ret, err = agentExec.Invoke(&req.BaseRequest)

	ret.Id = req.DebugInterfaceId

	return
}
