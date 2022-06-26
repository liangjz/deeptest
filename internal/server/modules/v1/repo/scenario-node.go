package repo

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type ScenarioNodeRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *ScenarioNodeRepo) GetTree(scenarioId int) (root *model.TestProcessor, err error) {
	pos, err := r.ListByScenario(scenarioId)
	if err != nil {
		return
	}

	root = pos[0]
	root.Slots = iris.Map{"icon": "icon"}
	r.makeTree(pos[1:], root)

	return
}

func (r *ScenarioNodeRepo) ListByScenario(scenarioId int) (pos []*model.TestProcessor, err error) {
	err = r.DB.
		Where("scenario_id=?", scenarioId).
		Where("NOT deleted").
		Order("parent_id ASC, ordr ASC").
		Find(&pos).Error
	return
}

func (r *ScenarioNodeRepo) makeTree(Data []*model.TestProcessor, node *model.TestProcessor) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := r.haveChild(Data, node) //判断节点是否有子节点并返回
	if children != nil {
		node.Children = append(node.Children, children[0:]...) //添加子节点
		for _, v := range children {                           //查询子节点的子节点，并添加到子节点
			_, has := r.haveChild(Data, v)
			if has {
				r.makeTree(Data, v) //递归添加节点
			}
		}
	}
}

func (r *ScenarioNodeRepo) haveChild(Data []*model.TestProcessor, node *model.TestProcessor) (child []*model.TestProcessor, yes bool) {
	for _, v := range Data {
		if v.ParentId == node.ID {
			v.Slots = iris.Map{"icon": "icon"}
			child = append(child, v)
		}
	}
	if child != nil {
		yes = true
	}
	return
}

func (r *ScenarioNodeRepo) CreateDefault(scenarioId uint) (po model.TestProcessor, err error) {
	po = model.TestProcessor{
		ScenarioId:     scenarioId,
		Name:           "root",
		EntityCategory: consts.ProcessorRoot,
		EntityId:       0,
	}
	err = r.DB.Create(&po).Error

	return
}

func (r *ScenarioNodeRepo) Save(processor *model.TestProcessor) (err error) {
	err = r.DB.Save(processor).Error

	return
}

//func (r *ScenarioNodeRepo) UpdateOrder(pos serverConsts.DropPos, targetId uint) (parentId uint, ordr int) {
//	if pos == serverConsts.Inner {
//		parentId = targetId
//
//		var preChild model.ScenarioNode
//		r.DB.Where("parent_id=?", parentId).
//			Order("ordr DESC").Limit(1).
//			First(&preChild)
//
//		ordr = preChild.Ordr + 1
//
//	} else if pos == serverConsts.Before {
//		brother, _ := r.Get(targetId)
//		parentId = brother.ParentId
//
//		r.DB.Model(&model.ScenarioNode{}).
//			Where("NOT deleted AND parent_id=? AND ordr >= ?", parentId, brother.Ordr).
//			Update("ordr", gorm.Expr("ordr + 1"))
//
//		ordr = brother.Ordr
//
//	} else if pos == serverConsts.After {
//		brother, _ := r.Get(targetId)
//		parentId = brother.ParentId
//
//		r.DB.Model(&model.ScenarioNode{}).
//			Where("NOT deleted AND parent_id=? AND ordr > ?", parentId, brother.Ordr).
//			Update("ordr", gorm.Expr("ordr + 1"))
//
//		ordr = brother.Ordr + 1
//
//	}
//
//	return
//}
//
//func (r *ScenarioNodeRepo) ListByScenario(projectId int) (pos []*model.ScenarioNode, err error) {
//	err = r.DB.
//		Where("project_id=?", projectId).
//		Where("NOT deleted").
//		Order("parent_id ASC, ordr ASC").
//		Find(&pos).Error
//	return
//}
//
//func (r *ScenarioNodeRepo) Get(fieldId uint) (field model.ScenarioNode, err error) {
//	err = r.DB.
//		Where("id=?", fieldId).
//		Where("NOT deleted").
//		First(&field).Error
//	return
//}
//
//func (r *ScenarioNodeRepo) SaveLogic(field *model.ScenarioNode) (err error) {
//	err = r.DB.SaveLogic(field).Error
//	return
//}
//
//func (r *ScenarioNodeRepo) Update(interf model.ScenarioNode) (err error) {
//	r.DB.Transaction(func(tx *gorm.DB) error {
//		err = r.DB.Updates(interf).Error
//		if err != nil {
//			return err
//		}
//
//		return err
//	})
//
//	return
//}

func (r *ScenarioNodeRepo) UpdateName(id int, name string) (err error) {
	err = r.DB.Model(&model.TestProcessor{}).
		Where("id = ?", id).
		Update("name", name).Error

	return
}

//func (r *ScenarioNodeRepo) makeTree(Data []*model.ScenarioNode, node *model.ScenarioNode) { //参数为父节点，添加父节点的子节点指针切片
//	children, _ := r.haveChild(Data, node) //判断节点是否有子节点并返回
//	if children != nil {
//		node.Children = append(node.Children, children[0:]...) //添加子节点
//		for _, v := range children {                           //查询子节点的子节点，并添加到子节点
//			_, has := r.haveChild(Data, v)
//			if has {
//				r.makeTree(Data, v) //递归添加节点
//			}
//		}
//	}
//}
//
//func (r *ScenarioNodeRepo) haveChild(Data []*model.ScenarioNode, node *model.ScenarioNode) (child []*model.ScenarioNode, yes bool) {
//	for _, v := range Data {
//		if v.ParentId == node.ID {
//			v.Slots = iris.Map{"icon": "icon"}
//			child = append(child, v)
//		}
//	}
//	if child != nil {
//		yes = true
//	}
//	return
//}
//
//func (r *ScenarioNodeRepo) Delete(id uint) (err error) {
//	err = r.DB.Model(&model.ScenarioNode{}).
//		Where("id=?", id).
//		Update("deleted", true).
//		Error
//
//	//field := model.ScenarioNode{}
//	//field.ID = id
//	//err = r.DB.Delete(field).Error
//
//	return
//}
//
//func (r *ScenarioNodeRepo) GetChildren(defId, fieldId uint) (children []*model.ScenarioNode, err error) {
//	err = r.DB.Where("defID=? AND parentID=?", defId, fieldId).Find(&children).Error
//	return
//}
//
//func (r *ScenarioNodeRepo) UpdateOrdAndParent(interf model.ScenarioNode) (err error) {
//	err = r.DB.Model(&interf).
//		Updates(model.ScenarioNode{Ordr: interf.Ordr, ParentId: interf.ParentId}).
//		Error
//
//	return
//}
