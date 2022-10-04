package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type ProjectModule struct {
	ProjectCtrl *handler.ProjectCtrl `inject:""`
}

func NewProjectModule() *ProjectModule {
	return &ProjectModule{}
}

// Party 项目
func (m *ProjectModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/", m.ProjectCtrl.List).Name = "项目列表"
		index.Get("/{id:uint}", m.ProjectCtrl.Get).Name = "项目详情"
		index.Post("/", m.ProjectCtrl.Create).Name = "新建项目"
		index.Put("/", m.ProjectCtrl.Update).Name = "更新项目"
		index.Delete("/{id:uint}", m.ProjectCtrl.Delete).Name = "删除项目"

		index.Post("/changeProject", m.ProjectCtrl.ChangeProject).Name = "切换用户默认项目"
		index.Get("/getByUser", m.ProjectCtrl.GetByUser).Name = "获取用户参与的项目"
	}
	return module.NewModule("/projects", handler)
}