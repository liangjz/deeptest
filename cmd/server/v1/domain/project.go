package domain

import (
	"github.com/aaronchen2k/deeptest/pkg/domain"
)

type ProjectReq struct {
	_domain.Model
	ProjectBase
}

type ProjectReqPaginate struct {
	_domain.PaginateReq
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type ProjectResp struct {
	_domain.PaginateReq
	ProjectBase
}

type ProjectMemberRemoveReq struct {
	UserId    int `json:"userId"`
	ProjectId int `json:"projectId"`
}

type ProjectBase struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc" gorm:"column:descr"`

	SchemaId       uint   `json:"schemaId"`
	OrgId          uint   `json:"orgId"`
	Logo           string `json:"logo"`
	ShortName      string `json:"shortName" validate:"required"`
	IncludeExample bool   `json:"includeExample"`
	AdminId        uint   `json:"adminId" validate:"required"`
	AdminName      string `gorm:"-" json:"adminName"`
}

type ProjectUserPermsPaginate struct {
	_domain.PaginateReq
}

type UpdateProjectMemberReq struct {
	ProjectId     uint `json:"projectId"`
	ProjectRoleId uint `json:"projectRoleId"`
	UserId        uint `json:"userId"`
}
