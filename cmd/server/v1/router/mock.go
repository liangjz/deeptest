package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/kataras/iris/v12"
)

type MockModule struct {
	MockCtrl *handler.MockCtrl `inject:""`
}

// Party 脚本
func (m *MockModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"
		index.Post("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"
		index.Put("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"
		index.Delete("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"
		index.Patch("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"
		index.Head("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"
		index.Connect("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"
		index.Trace("/{serveId:int}/{path:path}", m.MockCtrl.Mock).Name = "测试"

		index.Get("/", m.MockCtrl.Get).Name = "模拟接口测试"
		index.Post("/", m.MockCtrl.Posts).Name = "模拟接口测试"
		index.Put("/", m.MockCtrl.Posts).Name = "模拟接口测试"
		index.Delete("/", m.MockCtrl.Posts).Name = "模拟接口测试"

		index.Patch("/", m.MockCtrl.Posts).Name = "模拟接口测试"
		index.Head("/", m.MockCtrl.Head).Name = "模拟接口测试"

		index.Connect("/", m.MockCtrl.Connect).Name = "模拟接口测试"
		index.Trace("/", m.MockCtrl.Trace).Name = "模拟接口测试"
	}
	return module.NewModule("/", handler)
}
