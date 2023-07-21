package repo

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type ExtractorRepo struct {
	DB                *gorm.DB           `inject:""`
	PostConditionRepo *PostConditionRepo `inject:""`
}

func (r *ExtractorRepo) List(debugInterfaceId, endpointInterfaceId uint) (pos []model.DebugConditionExtractor, err error) {
	db := r.DB.
		Where("NOT deleted").
		Order("created_at ASC")

	if debugInterfaceId > 0 {
		db.Where("debug_interface_id=?", debugInterfaceId)
	} else {
		db.Where("endpoint_interface_id=? AND debug_interface_id=?", endpointInterfaceId, 0)
	}

	err = db.Find(&pos).Error

	return
}

//func (r *ExtractorRepo) ListTo(debugInterfaceId, endpointInterfaceId uint) (ret []domain.ExtractorBase, err error) {
//	pos, _ := r.List(debugInterfaceId, endpointInterfaceId)
//
//	for _, po := range pos {
//		extractor := domain.ExtractorBase{}
//		copier.CopyWithOption(&extractor, po, copier.Option{DeepCopy: true})
//
//		ret = append(ret, extractor)
//	}
//
//	return
//}

func (r *ExtractorRepo) Get(id uint) (extractor model.DebugConditionExtractor, err error) {
	err = r.DB.
		Where("id=?", id).
		Where("NOT deleted").
		First(&extractor).Error
	return
}

func (r *ExtractorRepo) GetByInterfaceVariable(variable string, id, debugInterfaceId uint) (extractor model.DebugConditionExtractor, err error) {
	db := r.DB.Model(&extractor).
		Where("variable = ? AND debug_interface_id =? AND not deleted",
			variable, debugInterfaceId)

	if id > 0 {
		db.Where("id != ?", id)
	}

	db.First(&extractor)

	return
}

func (r *ExtractorRepo) Save(extractor *model.DebugConditionExtractor) (id uint, bizErr _domain.BizErr) {
	//postCondition, _ := r.PostConditionRepo.Get(extractor.ConditionId)
	//
	//po, _ := r.GetByInterfaceVariable(extractor.Variable, extractor.ID, postCondition.EndpointInterfaceId)
	//if po.ID > 0 {
	//	bizErr.Code = _domain.ErrNameExist.Code
	//	return
	//}

	err := r.DB.Save(extractor).Error
	if err != nil {
		bizErr.Code = _domain.SystemErr.Code
		return
	}

	id = extractor.ID

	return
}

func (r *ExtractorRepo) Update(extractor *model.DebugConditionExtractor) (err error) {
	r.UpdateDesc(extractor)

	err = r.DB.Updates(extractor).Error
	if err != nil {
		return
	}

	return
}

func (r *ExtractorRepo) UpdateDesc(po *model.DebugConditionExtractor) (err error) {
	src := ""
	if po.Src == consts.Header {
		src = "响应头"
	} else if po.Src == consts.Body {
		src = "响应体"
	}

	name := ""
	if po.Type == consts.Boundary {
		name = fmt.Sprintf("边界提取器 \"%s - %s\"", po.BoundaryStart, po.BoundaryEnd)
	} else if po.Type == consts.JsonQuery {
		name = fmt.Sprintf("JSON提取器 \"%s\"", po.Expression)
	} else if po.Type == consts.HtmlQuery {
		name = fmt.Sprintf("HTML提取器 \"%s\"", po.Expression)
	} else if po.Type == consts.XmlQuery {
		name = fmt.Sprintf("XML提取器 \"%s\"", po.Expression)
	} else if po.Type == consts.Regx {
		name = fmt.Sprintf("正则表达式提取器 \"%s\"", po.Expression)
	}

	desc := fmt.Sprintf("%s%s", src, name)
	values := map[string]interface{}{
		"desc": desc,
	}

	err = r.DB.Model(&model.DebugPostCondition{}).
		Where("id=?", po.ConditionId).
		Updates(values).Error

	return
}

//func (r *ExtractorRepo) CreateOrUpdateResult(extractor *model.DebugConditionExtractor, usedBy consts.UsedBy) (err error) {
//	postCondition, _ := r.PostConditionRepo.Get(extractor.ConditionId)
//
//	po, _ := r.GetByInterfaceVariable(extractor.Variable, extractor.ID, postCondition.EndpointInterfaceId)
//	if po.ID > 0 {
//		extractor.ID = po.ID
//		r.UpdateResult(*extractor, usedBy)
//		return
//	}
//
//	err = r.DB.Save(extractor).Error
//	if err != nil {
//		return
//	}
//
//	return
//}

func (r *ExtractorRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.DebugConditionExtractor{}).
		Where("id=?", id).
		Update("deleted", true).
		Error

	return
}
func (r *ExtractorRepo) DeleteByCondition(conditionId uint) (err error) {
	err = r.DB.Model(&model.DebugConditionExtractor{}).
		Where("condition_id=?", conditionId).
		Update("deleted", true).
		Error

	return
}

func (r *ExtractorRepo) UpdateResult(extractor domain.ExtractorBase) (err error) {
	extractor.Result = strings.TrimSpace(extractor.Result)
	values := map[string]interface{}{}
	if extractor.Result != "" {
		values["result"] = extractor.Result
	}
	if extractor.Scope != "" {
		values["scope"] = extractor.Scope
	}

	err = r.DB.Model(&model.DebugConditionExtractor{}).
		Where("id = ?", extractor.RecordId).
		Updates(values).Error

	if err != nil {
		logUtils.Errorf("update DebugConditionExtractor error", zap.String("error:", err.Error()))
		return err
	}

	return
}

func (r *ExtractorRepo) CreateLog(extractor domain.ExtractorBase, invokeId uint) (
	logExtractor model.ExecLogExtractor, err error) {

	copier.CopyWithOption(&logExtractor, extractor, copier.Option{DeepCopy: true})

	logExtractor.ID = 0
	logExtractor.InvokeId = invokeId
	logExtractor.CreatedAt = nil
	logExtractor.UpdatedAt = nil

	err = r.DB.Save(&logExtractor).Error

	return
}

//func (r *ExtractorRepo) UpdateResultToExecLog(extractor model.DebugConditionExtractor, log *model.ExecLogProcessor) (
//	logExtractor model.ExecLogExtractor, err error) {
//
//	copier.CopyWithOption(&logExtractor, extractor, copier.Option{DeepCopy: true})
//
//	logExtractor.ID = 0
//	logExtractor.LogId = log.ID
//	logExtractor.CreatedAt = nil
//	logExtractor.UpdatedAt = nil
//
//	err = r.DB.Save(&logExtractor).Error
//
//	return
//}

func (r *ExtractorRepo) ListExtractorVariableByInterface(req domain.DebugReq) (variables []domain.Variable, err error) {
	err = r.DB.Model(&model.DebugConditionExtractor{}).
		Select("id, variable AS name, result AS value").
		Where("debug_interface_id=?", req.DebugInterfaceId).
		Where("NOT deleted AND NOT disabled").
		Order("created_at ASC").
		Find(&variables).Error

	return
}

func (r *ExtractorRepo) ListValidExtractorVariableForInterface(interfaceId, projectId uint, usedBy consts.UsedBy) (
	variables []domain.Variable, err error) {

	q := r.DB.Model(&model.DebugConditionExtractor{}).
		Select("id, variable AS name, result AS value, " +
			"endpoint_interface_id AS endpointInterfaceId, scope AS scope").
		Where("NOT deleted AND NOT disabled")

	//if usedBy == consts.InterfaceDebug {
	//	q.Where("project_id=?", projectId)
	//
	//} else {
	//	processorInterface, _ := r.ProcessorInterfaceRepo.Get(interfaceId)
	//
	//	var parentIds []uint
	//	r.GetParentIds(processorInterface.ProcessorId, &parentIds)
	//
	//	q.Where("scenario_id=?", processorInterface.ScenarioId).
	//		Where("scope = ? OR scenario_processor_id IN(?)", consts.Public, parentIds)
	//}

	err = q.Order("created_at ASC").
		Find(&variables).Error

	return
}

func (r *ExtractorRepo) GetParentIds(processorId uint, ids *[]uint) {
	var po model.Processor

	r.DB.Where("id = ?", processorId).
		Where("NOT deleted AND NOT disabled").
		First(&po)

	if po.ID > 0 {
		*ids = append(*ids, processorId)
	}

	if po.ParentId > 0 {
		r.GetParentIds(po.ParentId, ids)
	}

	return
}

func (r *ExtractorRepo) CloneFromEndpointInterfaceToDebugInterface(endpointInterfaceId, debugInterfaceId uint,
	usedBy consts.UsedBy) (
	err error) {

	srcPos, _ := r.List(0, endpointInterfaceId)

	for _, po := range srcPos {
		po.ID = 0
		//po.EndpointInterfaceId = endpointInterfaceId
		//po.DebugInterfaceId = debugInterfaceId
		//po.UsedBy = usedBy

		r.Save(&po)
	}

	return
}

func (r *ExtractorRepo) CreateDefault(conditionId uint) (po model.DebugConditionExtractor) {
	po = model.DebugConditionExtractor{
		ConditionId: conditionId,

		ExtractorBase: domain.ExtractorBase{
			Src:        consts.Body,
			Type:       consts.Boundary,
			Expression: "",
			Variable:   "",
			Scope:      consts.Public,
		},
	}

	r.Save(&po)

	return
}
