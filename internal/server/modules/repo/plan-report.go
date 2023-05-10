package repo

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	_commUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"gorm.io/gorm"
)

type PlanReportRepo struct {
	DB       *gorm.DB  `inject:""`
	LogRepo  *LogRepo  `inject:""`
	UserRepo *UserRepo `inject:""`
}

func (r *PlanReportRepo) Paginate(req v1.ReportReqPaginate, projectId int) (data _domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.ScenarioReport{}).
		Where("project_id = ? AND NOT deleted", projectId)

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.ScenarioId != 0 {
		db = db.Where("scenario_id = ?", req.ScenarioId)
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count report error %s", err.Error())
		return
	}

	results := make([]*model.ScenarioReport, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&results).Error
	if err != nil {
		logUtils.Errorf("query report error %s", err.Error())
		return
	}

	r.CombineUserName(results)
	data.Populate(results, count, req.Page, req.PageSize)

	return
}

func (r *PlanReportRepo) CombineUserName(data []*model.ScenarioReport) {
	userIds := make([]uint, 0)
	for _, v := range data {
		userIds = append(userIds, v.CreateUserId)
	}
	userIds = _commUtils.ArrayRemoveUintDuplication(userIds)

	users, _ := r.UserRepo.FindByIds(userIds)

	userIdNameMap := make(map[uint]string)
	for _, v := range users {
		userIdNameMap[v.ID] = v.Name
	}

	for _, v := range data {
		if name, ok := userIdNameMap[v.CreateUserId]; ok {
			v.CreateUserName = name
		}
	}
}

func (r *PlanReportRepo) Get(id uint) (report model.ScenarioReport, err error) {
	err = r.DB.Where("id = ?", id).First(&report).Error
	if err != nil {
		logUtils.Errorf("find report by id error %s", err.Error())
		return
	}

	root, err := r.getLogTree(report)
	report.Logs = root.Logs

	return
}

func (r *PlanReportRepo) Create(result *model.PlanReport) (bizErr *_domain.BizErr) {
	err := r.DB.Model(&model.PlanReport{}).Create(result).Error
	if err != nil {
		logUtils.Errorf("create plan report error %s", err.Error())
		bizErr.Code = _domain.SystemErr.Code

		return
	}

	return
}

func (r *PlanReportRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.PlanReport{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete report by id error %s", err.Error())
		return
	}

	err = r.DB.Model(&model.ExecLogProcessor{}).Where("report_id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete report's logs by id error %s", err.Error())
		return
	}

	return
}

func (r *PlanReportRepo) UpdateStatus(progressStatus consts.ProgressStatus, resultStatus consts.ResultStatus, scenarioId uint) (
	err error) {

	values := map[string]interface{}{
		"progress_status": progressStatus,
		"result_status":   resultStatus,
	}
	err = r.DB.Model(&model.ScenarioReport{}).
		Where("report_id = ? AND progress_status = ?", scenarioId, consts.InProgress).
		Updates(values).Error

	return
}

func (r *PlanReportRepo) UpdateResult(report model.ScenarioReport) (err error) {
	values := map[string]interface{}{
		"pass_num":        report.PassRequestNum,
		"fail_num":        report.FailRequestNum,
		"start_time":      report.StartTime,
		"end_time":        report.EndTime,
		"duration":        report.Duration,
		"progress_status": consts.End,
		"result_status":   report.ResultStatus,
	}
	err = r.DB.Model(&model.ScenarioReport{}).
		Where("id = ?", report.ID).
		Updates(values).Error

	return
}

func (r *PlanReportRepo) ResetResult(result model.ScenarioReport) (err error) {
	values := map[string]interface{}{
		"name":       result.Name,
		"start_time": result.StartTime,
	}
	err = r.DB.Model(&result).Where("id = ?", result.ID).Updates(values).Error
	if err != nil {
		logUtils.Errorf("update report error %s", err.Error())
		return
	}

	return
}

func (r *PlanReportRepo) ClearLogs(resultId uint) (err error) {
	err = r.DB.Model(&model.ExecLogProcessor{}).Where("result_id = ?", resultId).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete logs by result id error %s", err.Error())
		return
	}

	return
}

func (r *PlanReportRepo) FindInProgressResult(scenarioId uint) (result model.ScenarioReport, err error) {
	err = r.DB.Model(&result).
		Where("progress_status =? AND scenario_id = ? AND  not deleted", consts.InProgress, scenarioId).
		First(&result).Error

	return
}

func (r *PlanReportRepo) getLogTree(report model.ScenarioReport) (root model.ExecLogProcessor, err error) {
	logs, err := r.LogRepo.ListByReport(report.ID)
	if err != nil {
		return
	}

	for _, log := range logs {
		if log.ProcessorType == consts.ProcessorInterfaceDefault {
			log.InterfaceExtractorsResult, _ = r.listLogExtractors(log.ID)
			log.InterfaceCheckpointsResult, _ = r.listLogCheckpoints(log.ID)
		}
	}

	root = model.ExecLogProcessor{
		Name: report.Name,
	}
	r.makeTree(logs, &root)

	return
}

func (r *PlanReportRepo) makeTree(Data []*model.ExecLogProcessor, parent *model.ExecLogProcessor) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := r.haveChild(Data, parent) //判断节点是否有子节点并返回

	if children != nil {
		parent.Logs = append(parent.Logs, children[0:]...) //添加子节点

		for _, child := range children { //查询子节点的子节点，并添加到子节点
			_, has := r.haveChild(Data, child)
			if has {
				r.makeTree(Data, child) //递归添加节点
			}
		}
	}
}

func (r *PlanReportRepo) haveChild(Data []*model.ExecLogProcessor, node *model.ExecLogProcessor) (children []*model.ExecLogProcessor, yes bool) {
	for _, v := range Data {
		if v.ParentId == node.ID {
			children = append(children, v)
		}
	}

	if children != nil {
		yes = true
	}

	return
}

func (r *PlanReportRepo) listLogExtractors(logId uint) (extractors []model.ExecLogExtractor, err error) {
	err = r.DB.
		Where("log_id =? AND not deleted", logId).
		Find(&extractors).Error

	return
}

func (r *PlanReportRepo) listLogCheckpoints(logId uint) (checkpoints []model.ExecLogCheckpoint, err error) {
	err = r.DB.
		Where("log_id =? AND not deleted", logId).
		Find(&checkpoints).Error

	return
}

func (r *PlanReportRepo) GetLastByPlanId(planId uint) (report model.PlanReport, err error) {
	err = r.DB.Where("plan_id = ?", planId).Last(&report).Error
	if err != nil {
		logUtils.Errorf("find plan report by plan_id error %s", err.Error())
		return
	}
	return
}
