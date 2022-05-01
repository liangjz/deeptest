package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type InvocationModule struct {
	InvocationCtrl *controller.InvocationCtrl `inject:""`
}

// Party 脚本
func (m *InvocationModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/", m.InvocationCtrl.List).Name = "调用列表"
		index.Get("/{id:uint}", m.InvocationCtrl.GetAsInterface).Name = "调用详情"
		index.Delete("/{id:uint}", m.InvocationCtrl.Delete).Name = "删除调用"
	}
	return module.NewModule("/invocations", handler)
}