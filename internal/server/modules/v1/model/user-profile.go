package model

type SysUserProfile struct {
	BaseModel

	Email         string `json:"email"`
	Phone         string `json:"phone"`
	CurrProjectId uint   `json:"currProjectId"`

	UserId string `json:"userId"`
}

func (SysUserProfile) TableName() string {
	return "sys_user_profile"
}