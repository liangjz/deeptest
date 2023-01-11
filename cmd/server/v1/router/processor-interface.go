package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type ProcessorInterfaceModule struct {
	ProcessorInterfaceCtrl *handler.ProcessorInterfaceCtrl `inject:""`
}

// Party 场景
func (m *ProcessorInterfaceModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/getInterface", m.ProcessorInterfaceCtrl.GetInterface).Name = "获取场景接口处理器"
		index.Get("/listInvocation", m.ProcessorInterfaceCtrl.ListInvocation).Name = "获取场景接口调用历史"
	}

	return module.NewModule("/processors/interfaces", handler)
}
