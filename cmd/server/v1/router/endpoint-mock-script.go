package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type EndpointMockScriptModule struct {
	EndpointMockScriptCtrl *handler.EndpointMockScriptCtrl `inject:""`
}

// Party 项目
func (m *EndpointMockScriptModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/{id:uint}", m.EndpointMockScriptCtrl.Get).Name = "项目详情"
		index.Put("/", m.EndpointMockScriptCtrl.Update).Name = "更新项目"
	}
	return module.NewModule("/mockScripts", handler)
}
