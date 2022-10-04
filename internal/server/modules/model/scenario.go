package model

type Scenario struct {
	BaseModel

	Version float64 `json:"version" yaml:"version"`
	Name    string  `json:"name" yaml:"name"`
	Desc    string  `json:"desc" yaml:"desc"`

	Processor Processor `json:"processor" yaml:"processor" gorm:"-"`

	ProjectId uint `json:"projectId"`
}

func (Scenario) TableName() string {
	return "biz_scenario"
}