package model

type Datapool struct {
	BaseModel

	Name string `json:"name"`
	Desc string `json:"desc"`
	Data string `json:"data"`
	Path string `json:"path"`

	ProjectId uint `json:"projectId"`
}

func (Datapool) TableName() string {
	return "biz_datapool"
}