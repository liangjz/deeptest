package handler

import (
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type ParserCtrl struct {
	ParserService *service.ParserService `inject:""`
	BaseCtrl
}

// ParseHtml
func (c *ParserCtrl) ParseHtml(ctx iris.Context) {
	req := v1.ParserRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	resp, err := c.ParserService.ParseHtml(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp, Msg: _domain.NoErr.Msg})
}

// TestXPath
func (c *ParserCtrl) TestXPath(ctx iris.Context) {
	req := v1.TestXPathRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	resp, err := c.ParserService.TestXPath(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp, Msg: _domain.NoErr.Msg})
}