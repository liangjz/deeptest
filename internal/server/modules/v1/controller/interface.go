package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type InterfaceCtrl struct {
	InterfaceService  *service.InterfaceService  `inject:""`
	InvocationService *service.InvocationService `inject:""`
	ExtractorService  *service.ExtractorService  `inject:""`

	BaseCtrl
}

// InvokeInterface
func (c *InterfaceCtrl) InvokeInterface(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("currProjectId")
	if projectId == 0 {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: "projectId"})
		return
	}

	req := serverDomain.InvocationRequest{}
	err = ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.InterfaceService.UpdateByRequest(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	resp, err := c.InterfaceService.Test(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	_, err = c.InvocationService.Create(req, resp, projectId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	c.ExtractorService.ExtractByInterface(req.Id, resp)

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp})
}

// SaveInterface
func (c *InterfaceCtrl) SaveInterface(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("currProjectId")
	if projectId == 0 {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: "projectId"})
		return
	}

	req := serverDomain.InvocationRequest{}
	err = ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.InterfaceService.UpdateByConfig(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	resp, err := c.InterfaceService.Test(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp})
}

// Load
func (c *InterfaceCtrl) Load(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("currProjectId")

	data, err := c.InterfaceService.GetTree(projectId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data, Msg: _domain.NoErr.Msg})
}

// Get 详情
func (c *InterfaceCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Data: nil, Msg: _domain.ParamErr.Msg})
		return
	}

	interf, err := c.InterfaceService.Get(id)

	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: interf})
}

// Create 添加
func (c *InterfaceCtrl) Create(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("currProjectId")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	req := serverDomain.InterfaceReq{}
	err = ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	req.ProjectId = projectId
	intf, err := c.InterfaceService.Create(req)
	if err != nil {
		ctx.JSON(_domain.Response{
			Code: c.ErrCode(err),
			Data: nil,
		})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: intf, Msg: _domain.NoErr.Msg})
}

// Update 更新
func (c *InterfaceCtrl) Update(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")

	var req serverDomain.InterfaceReq
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err = c.InterfaceService.Update(id, req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: nil, Msg: _domain.NoErr.Msg})
}

// UpdateName 更新
func (c *InterfaceCtrl) UpdateName(ctx iris.Context) {
	var req serverDomain.InterfaceReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		logUtils.Errorf("参数验证失败", err.Error())
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	err = c.InterfaceService.UpdateName(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: nil, Msg: _domain.NoErr.Msg})
}

// Delete 删除
func (c *InterfaceCtrl) Delete(ctx iris.Context) {
	projectId, _ := ctx.URLParamInt("currProjectId")

	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.InterfaceService.Delete(uint(projectId), uint(id))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

func (c *InterfaceCtrl) Move(ctx iris.Context) {
	projectId, _ := ctx.URLParamInt("currProjectId")

	var req serverDomain.InterfaceMoveReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	_, err = c.InterfaceService.Move(uint(req.DragKey), uint(req.DropKey), req.DropPos, uint(projectId))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}