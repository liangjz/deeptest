package model

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type ComponentSchema struct {
	BaseModel
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Content     string            `json:"content" gorm:"type:text"`
	ServeId     int64             `json:"serveId"`
	Examples    string            `json:"examples" gorm:"type:text"`
	Tags        string            `json:"tags"`
	Description string            `json:"description"`
	Ref         string            `json:"ref"`
	SourceType  consts.SourceType `json:"sourceType" gorm:"default:''"`
}

func (ComponentSchema) TableName() string {
	return "biz_project_serve_component_schema"
}
