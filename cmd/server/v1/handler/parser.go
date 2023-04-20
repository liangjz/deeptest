package handler

import (
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type ParserCtrl struct {
	ParserService     *service.ParserService     `inject:""`
	ParserHtmlService *service.ParserHtmlService `inject:""`
	ParserXmlService  *service.ParserXmlService  `inject:""`
	ParserJsonService *service.ParserJsonService `inject:""`
	ParserRegxService *service.ParserRegxService `inject:""`
	BaseCtrl
}

// ParseHtml
func (c *ParserCtrl) ParseHtml(ctx iris.Context) {
	req := v1.ParserRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	resp, err := c.ParserHtmlService.ParseHtml(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp, Msg: _domain.NoErr.Msg})
}

// ParseXml
func (c *ParserCtrl) ParseXml(ctx iris.Context) {
	req := v1.ParserRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	resp, err := c.ParserXmlService.ParseXml(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp, Msg: _domain.NoErr.Msg})
}

// ParseJson
func (c *ParserCtrl) ParseJson(ctx iris.Context) {
	req := v1.ParserRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	resp, err := c.ParserJsonService.ParseJson(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp, Msg: _domain.NoErr.Msg})
}

// ParseText
func (c *ParserCtrl) ParseText(ctx iris.Context) {
	req := v1.ParserRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	resp, err := c.ParserRegxService.ParseRegx(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp, Msg: _domain.NoErr.Msg})
}

// TestExpr
func (c *ParserCtrl) TestExpr(ctx iris.Context) {
	req := v1.TestExprRequest{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	resp, err := c.ParserService.TestExpr(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.FailErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp, Msg: _domain.NoErr.Msg})
}
