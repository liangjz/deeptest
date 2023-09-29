package repo

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"gorm.io/gorm"
)

type JslibRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *JslibRepo) List(keywords string) (pos []model.SysJslib, err error) {
	db := r.DB.Model(&model.SysJslib{}).
		Where("NOT deleted")

	if keywords != "" {
		db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", keywords))
	}

	err = db.Find(&pos).Error

	return
}

func (r *JslibRepo) Get(id uint) (po model.SysJslib, err error) {
	err = r.DB.Model(&model.SysJslib{}).
		Where("id = ?", id).First(&po).Error

	return
}

func (r *JslibRepo) Save(po *model.SysJslib) (err error) {
	err = r.DB.Save(po).Error

	return
}

func (r *JslibRepo) UpdateName(to v1.JslibReq) (err error) {
	err = r.DB.Model(&model.SysJslib{}).
		Where("id = ?", to.Id).
		Updates(map[string]interface{}{"name": to.Name, "update_user": to.UpdateUser}).Error

	return
}

func (r *JslibRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.SysJslib{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error

	return
}

func (r *JslibRepo) Disable(id uint) (err error) {
	err = r.DB.Model(&model.SysJslib{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"disabled": gorm.Expr("NOT disabled")}).Error

	return
}