package handler

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	service "github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type DebugInvokeCtrl struct {
	DebugInvokeService *service.DebugInvokeService `inject:""`
	ExtractorService   *service.ExtractorService   `inject:""`
	CheckpointService  *service.CheckpointService  `inject:""`
	BaseCtrl
}

// SubmitResult
func (c *DebugInvokeCtrl) SubmitResult(ctx iris.Context) {
	req := domain.SubmitDebugResultRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	err = c.DebugInvokeService.SubmitResult(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code})
}

// List
func (c *DebugInvokeCtrl) List(ctx iris.Context) {
	endpointInterfaceId, err := ctx.URLParamInt("endpointInterfaceId")
	debugInterfaceId, err := ctx.URLParamInt("debugInterfaceId")

	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	data, err := c.DebugInvokeService.ListByInterface(endpointInterfaceId, debugInterfaceId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data})
}

// GetLastResp
func (c *DebugInvokeCtrl) GetLastResp(ctx iris.Context) {
	endpointInterfaceId, err := ctx.URLParamInt("endpointInterfaceId")
	debugInterfaceId, err := ctx.URLParamInt("debugInterfaceId")

	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	resp, err := c.DebugInvokeService.GetLastResp(endpointInterfaceId, debugInterfaceId)

	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp})
}

// Get 详情
func (c *DebugInvokeCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	req, resp, err := c.DebugInvokeService.GetAsInterface(id)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: iris.Map{"req": req, "resp": resp}})
}

// Delete 删除
func (c *DebugInvokeCtrl) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.DebugInvokeService.Delete(uint(id))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}