package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"time"
)

type PlanReqPaginate struct {
	_domain.PaginateReq

	CategoryId int64             `json:"categoryId"`
	Status     consts.TestStatus `json:"status"`
	DirectorId uint              `json:"directorId"`
	Keywords   string            `json:"keywords"`
	Enabled    string            `json:"enabled"`
}

//type PlanAddScenariosReq struct {
//	SelectedNodes []ScenarioSimple `json:"selectedNodes"`
//
//	TargetId  uint `json:"targetId"`
//	ProjectId int  `json:"projectId"`
//}

type PlanAddScenariosReq struct {
	ScenarioIds []int `json:"scenarioIds"`
}

type PlanAndReportDetail struct {
	Id             uint              `json:"id"`           //计划ID
	AdminName      string            `json:"directorName"` //负责人姓名
	CreatedAt      *time.Time        `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time        `json:"updatedAt,omitempty"`
	UpdateUserName string            `json:"updaterName"`        //最近更新人姓名
	Status         consts.TestStatus `json:"status"`             //状态
	TestPassRate   string            `json:"testPassRate"`       //执行通过率
	ExecTimes      int64             `json:"execTimes"`          //执行次数
	ExecutorName   string            `json:"executorName"`       //执行人姓名
	ExecTime       *time.Time        `json:"execTime,omitempty"` //执行时间
	ExecEnv        string            `json:"execEnv"`            //执行环境
}
