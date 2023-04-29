package handler

import (
	"github.com/aaronchen2k/deeptest/cmd/agent/v1/domain"
	service "github.com/aaronchen2k/deeptest/internal/agent/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type ExecInterfaceCtrl struct {
	ExecService *service.ExecInterfaceService `inject:""`
}

// Call
func (c *ExecInterfaceCtrl) Call(ctx iris.Context) {
	req := domain.InterfaceCall{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	resp, err := c.ExecService.Run(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: resp})
}
