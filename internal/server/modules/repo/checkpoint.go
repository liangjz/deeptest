package repo

import (
	model2 "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type CheckpointRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *CheckpointRepo) List(interfaceId uint) (pos []model2.InterfaceCheckpoint, err error) {
	err = r.DB.
		Where("interface_id=?", interfaceId).
		Where("NOT deleted").
		Order("created_at ASC").
		Find(&pos).Error
	return
}

func (r *CheckpointRepo) Get(id uint) (checkpoint model2.InterfaceCheckpoint, err error) {
	err = r.DB.
		Where("id=?", id).
		Where("NOT deleted").
		First(&checkpoint).Error
	return
}

func (r *CheckpointRepo) GetByName(name string, interfaceId uint) (checkpoint model2.InterfaceCheckpoint, err error) {
	var checkpoints []model2.InterfaceCheckpoint

	db := r.DB.Model(&checkpoint).
		Where("name = ? AND interface_id =? AND not deleted", name, interfaceId)

	err = db.Find(&checkpoints).Error

	if err != nil {
		return
	}

	if len(checkpoints) > 0 {
		checkpoint = checkpoints[0]
	}

	return
}

func (r *CheckpointRepo) Save(checkpoint *model2.InterfaceCheckpoint) (err error) {
	err = r.DB.Save(checkpoint).Error
	return
}

func (r *CheckpointRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model2.InterfaceCheckpoint{}).
		Where("id=?", id).
		Update("deleted", true).
		Error

	return
}

func (r *CheckpointRepo) UpdateResult(checkpoint model2.InterfaceCheckpoint) (err error) {
	err = r.DB.Model(&checkpoint).
		Where("id=?", checkpoint.ID).
		Update("actual_result", checkpoint.ActualResult).
		Update("result_status", checkpoint.ResultStatus).
		Error

	return
}

func (r *CheckpointRepo) UpdateResultToExecLog(checkpoint model2.InterfaceCheckpoint, log *model2.ExecLogProcessor) (
	logCheckpoint model2.ExecLogCheckpoint, err error) {

	copier.CopyWithOption(&logCheckpoint, checkpoint, copier.Option{DeepCopy: true})
	logCheckpoint.ID = 0
	logCheckpoint.LogId = log.ID
	logCheckpoint.CreatedAt = nil
	logCheckpoint.UpdatedAt = nil

	err = r.DB.Save(&logCheckpoint).Error

	return
}