package handler

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type ExtractorCtrl struct {
	ExtractorService *service.ExtractorService `inject:""`
	BaseCtrl
}

// List
func (c *ExtractorCtrl) List(ctx iris.Context) {
	interfaceId, err := ctx.URLParamInt("interfaceId")
	if interfaceId == 0 {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: "interfaceId"})
		return
	}

	data, err := c.ExtractorService.List(interfaceId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data})
}

// Get 详情
func (c *ExtractorCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Data: nil, Msg: _domain.ParamErr.Msg})
		return
	}

	extractor, err := c.ExtractorService.Get(uint(id))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: extractor})
}

// Create 添加
func (c *ExtractorCtrl) Create(ctx iris.Context) {
	extractor := model.InterfaceExtractor{}
	err := ctx.ReadJSON(&extractor)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	bizErr := c.ExtractorService.Create(&extractor)
	if err != nil {
		ctx.JSON(_domain.Response{
			Code: bizErr.Code,
			Data: nil,
		})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: extractor, Msg: _domain.NoErr.Msg})
}

// Update 更新
func (c *ExtractorCtrl) Update(ctx iris.Context) {
	var extractor model.InterfaceExtractor
	err := ctx.ReadJSON(&extractor)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.ExtractorService.Update(&extractor)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: nil, Msg: _domain.NoErr.Msg})
}

// SaveOrUpdateResult 新建或更新结果
func (c *ExtractorCtrl) SaveOrUpdateResult(ctx iris.Context) {
	var extractor model.InterfaceExtractor
	err := ctx.ReadJSON(&extractor)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.ExtractorService.SaveOrUpdateResult(&extractor)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: nil, Msg: _domain.NoErr.Msg})
}

// Delete 删除
func (c *ExtractorCtrl) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.ExtractorService.Delete(uint(id))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// ListExtractorVariable
func (c *ExtractorCtrl) ListExtractorVariable(ctx iris.Context) {
	interfaceId, err := ctx.URLParamInt("interfaceId")
	if interfaceId == 0 {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: "interfaceId"})
		return
	}

	data, err := c.ExtractorService.ListExtractorVariableByInterface(interfaceId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data})
}

// ListValidExtractorVariable
func (c *ExtractorCtrl) ListValidExtractorVariable(ctx iris.Context) {
	interfaceId, err := ctx.URLParamInt("interfaceId")
	if interfaceId == 0 {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: "interfaceId"})
		return
	}

	data, err := c.ExtractorService.ListValidExtractorVariable(interfaceId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data})
}
