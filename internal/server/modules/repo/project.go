package repo

import (
	"errors"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	_commUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	commonUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	_fileUtils "github.com/aaronchen2k/deeptest/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	DB                         *gorm.DB                    `inject:""`
	RoleRepo                   *RoleRepo                   `inject:""`
	ProjectRoleRepo            *ProjectRoleRepo            `inject:""`
	EnvironmentRepo            *EnvironmentRepo            `inject:""`
	UserRepo                   *UserRepo                   `inject:""`
	ServeRepo                  *ServeRepo                  `inject:""`
	EndpointRepo               *EndpointRepo               `inject:""`
	ProjectRecentlyVisitedRepo *ProjectRecentlyVisitedRepo `inject:""`
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{}
}

func (r *ProjectRepo) Paginate(req v1.ProjectReqPaginate, userId uint) (data _domain.PageData, err error) {
	var count int64

	var projectIds []uint
	r.DB.Model(&model.ProjectMember{}).
		Select("project_id").Where("user_id = ?", userId).Scan(&projectIds)

	db := r.DB.Model(&model.Project{}).Where("NOT deleted AND id IN (?)", projectIds)

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", commonUtils.IsDisable(req.Enabled))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count project error", zap.String("error:", err.Error()))
		return
	}

	projects := make([]*model.Project, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&projects).Error
	if err != nil {
		logUtils.Errorf("query project error", zap.String("error:", err.Error()))
		return
	}

	for key, project := range projects {
		user, _ := r.UserRepo.FindById(project.AdminId)
		projects[key].AdminName = user.Name
	}

	data.Populate(projects, count, req.Page, req.PageSize)

	return
}

func (r *ProjectRepo) Get(id uint) (project model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("id = ?", id).
		First(&project).Error

	return
}

func (r *ProjectRepo) GetByName(projectName string, id uint) (project model.Project, err error) {
	db := r.DB.Model(&model.Project{}).
		Where("name = ?", projectName)

	if id > 0 {
		db.Where("id != ?", id)
	}

	err = db.First(&project).Error

	return
}

func (r *ProjectRepo) GetBySpec(spec string) (project model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("spec = ?", spec).
		First(&project).Error

	return
}

func (r *ProjectRepo) Save(po *model.Project) (err error) {
	err = r.DB.Save(po).Error

	return
}

func (r *ProjectRepo) Create(req v1.ProjectReq, userId uint) (id uint, bizErr *_domain.BizErr) {
	bizErr = new(_domain.BizErr)
	po, err := r.GetByName(req.Name, 0)
	if po.Name != "" {
		bizErr.Code = _domain.ErrNameExist.Code
		return
	}

	// create project
	project := model.Project{ProjectBase: req.ProjectBase}
	err = r.DB.Model(&model.Project{}).Create(&project).Error
	if err != nil {
		logUtils.Errorf("add project error", zap.String("error:", err.Error()))
		bizErr.Code = _domain.SystemErr.Code

		return
	}
	err = r.AddProjectMember(project.ID, req.AdminId, "admin")
	if err != nil {
		logUtils.Errorf("添加项目角色错误", zap.String("错误:", err.Error()))
		return
	}
	err = r.CreateProjectRes(project.ID, userId, req.IncludeExample)

	id = project.ID

	return
}

func (r *ProjectRepo) CreateProjectRes(projectId, userId uint, IncludeExample bool) (err error) {

	// create project member
	err = r.AddProjectMember(projectId, userId, "admin")
	if err != nil {
		logUtils.Errorf("添加项目角色错误", zap.String("错误:", err.Error()))
		return
	}

	// create project environment
	err = r.EnvironmentRepo.AddDefaultForProject(projectId)
	if err != nil {
		logUtils.Errorf("添加项目默认环境错误", zap.String("错误:", err.Error()))
		return
	}

	// create project serve
	serve, err := r.AddProjectDefaultServe(projectId, userId)
	if err != nil {
		logUtils.Errorf("添加默认服务错误", zap.String("错误:", err.Error()))
		return
	}

	// create project endpoint category
	err = r.AddProjectRootEndpointCategory(serve.ID, projectId)
	if err != nil {
		logUtils.Errorf("添加终端分类错误", zap.String("错误:", err.Error()))
		return
	}

	// create project scenario category
	err = r.AddProjectRootScenarioCategory(projectId)
	if err != nil {
		logUtils.Errorf("添加场景分类错误", zap.String("错误:", err.Error()))
		return
	}

	// create project plan category
	err = r.AddProjectRootPlanCategory(projectId)
	if err != nil {
		logUtils.Errorf("添加场景分类错误", zap.String("错误:", err.Error()))
		return
	}

	//create sample
	if IncludeExample {
		err = r.CreateSample(projectId, serve.ID, userId)
		if err != nil {
			logUtils.Errorf("创建示例失败", zap.String("错误:", err.Error()))
			return
		}
	}

	return
}

func (r *ProjectRepo) Update(req v1.ProjectReq) error {
	project := model.Project{ProjectBase: req.ProjectBase}
	err := r.DB.Model(&model.Project{}).Where("id = ?", req.Id).Updates(&project).Error
	if err != nil {
		logUtils.Errorf("update project error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *ProjectRepo) UpdateDefaultEnvironment(projectId, envId uint) (err error) {
	err = r.DB.Model(&model.Project{}).
		Where("id = ?", projectId).
		Updates(map[string]interface{}{"environment_id": envId}).Error

	if err != nil {
		logUtils.Errorf("update project environment error", err.Error())
		return err
	}

	return
}

func (r *ProjectRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.Project{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete project by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) DeleteChildren(ids []int, tx *gorm.DB) (err error) {
	err = tx.Model(&model.Project{}).Where("id IN (?)", ids).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("batch delete project error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *ProjectRepo) GetChildrenIds(id uint) (ids []int, err error) {
	tmpl := `
		WITH RECURSIVE project AS (
			SELECT * FROM biz_project WHERE id = %d
			UNION ALL
			SELECT child.* FROM biz_project child, project WHERE child.parent_id = project.id
		)
		SELECT id FROM project WHERE id != %d
    `
	sql := fmt.Sprintf(tmpl, id, id)
	err = r.DB.Raw(sql).Scan(&ids).Error
	if err != nil {
		logUtils.Errorf("get children project error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) ListProjectByUser(userId uint) (projects []model.ProjectMemberRole, err error) {
	isAdminUser, err := r.UserRepo.IsAdminUser(userId)
	if err != nil {
		return
	}

	db := r.DB.Model(&model.ProjectMember{}).
		Joins("LEFT JOIN biz_project p ON biz_project_member.project_id=p.id").
		Joins("LEFT JOIN biz_project_role r ON biz_project_member.project_role_id=r.id").
		Select("p.*, r.id role_id, r.name role_name").
		Where("NOT biz_project_member.deleted")

	if !isAdminUser {
		db.Where("biz_project_member.user_id = ?", userId)
	}
	err = db.Group("biz_project_member.project_id").Find(&projects).Error
	return
}

func (r *ProjectRepo) GetCurrProjectByUser(userId uint) (currProject model.Project, err error) {
	var user model.SysUser
	err = r.DB.Preload("Profile").
		Where("id = ?", userId).
		First(&user).
		Error

	err = r.DB.Model(&model.Project{}).
		Where("id = ?", user.Profile.CurrProjectId).
		First(&currProject).Error

	return
}

func (r *ProjectRepo) ListProjectsRecentlyVisited(userId uint) (projects []model.Project, err error) {
	err = r.DB.Raw(fmt.Sprintf("SELECT p.*,max( v.created_at ) visited_time FROM biz_project_recently_visited v,biz_project p WHERE v.project_id = p.id AND v.user_id = %d GROUP BY v.project_id ORDER BY visited_time DESC LIMIT 3", userId)).Find(&projects).Error
	return
}

func (r *ProjectRepo) ChangeProject(projectId, userId uint) (err error) {
	err = r.DB.Model(&model.SysUserProfile{}).Where("user_id = ?", userId).
		Updates(map[string]interface{}{"curr_project_id": projectId}).Error

	return
}

func (r *ProjectRepo) AddProjectMember(projectId, userId uint, role consts.RoleType) (err error) {
	var projectRole model.ProjectRole
	projectRole, err = r.ProjectRoleRepo.FindByName(role)
	if err != nil {
		return
	}

	projectMember := model.ProjectMember{UserId: userId, ProjectId: projectId, ProjectRoleId: projectRole.ID}
	err = r.DB.Create(&projectMember).Error

	return
}

func (r *ProjectRepo) AddProjectRootEndpointCategory(serveId, projectId uint) (err error) {
	root := model.Category{
		Name:      "分类",
		Type:      serverConsts.EndpointCategory,
		ServeId:   serveId,
		ProjectId: projectId,
		IsLeaf:    false,
	}
	err = r.DB.Create(&root).Error

	return
}

func (r *ProjectRepo) AddProjectRootScenarioCategory(projectId uint) (err error) {
	root := model.Category{
		Name:      "分类",
		Type:      serverConsts.ScenarioCategory,
		ProjectId: projectId,
		IsLeaf:    false,
	}
	err = r.DB.Create(&root).Error

	return
}

func (r *ProjectRepo) AddProjectRootPlanCategory(projectId uint) (err error) {
	root := model.Category{
		Name:      "分类",
		Type:      serverConsts.PlanCategory,
		ProjectId: projectId,
		IsLeaf:    false,
	}
	err = r.DB.Create(&root).Error

	return
}

func (r *ProjectRepo) Members(req v1.ProjectReqPaginate, projectId int) (data _domain.PageData, err error) {
	req.Order = "sys_user.created_at"
	db := r.DB.Model(&model.SysUser{}).
		Select("sys_user.id, sys_user.username, sys_user.email,sys_user.name, m.project_role_id, r.name as role_name").
		Joins("left join biz_project_member m on sys_user.id=m.user_id").
		Joins("left join biz_project_role r on m.project_role_id=r.id").
		Where("m.project_id = ?", projectId)
	if req.Keywords != "" {
		db = db.Where("sys_user.username LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}

	var count int64
	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count users error", zap.String("error:", err.Error()))
		return
	}

	users := make([]v1.MemberResp, 0)
	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, "", req.Order)).
		Scan(&users).Error
	if err != nil {
		logUtils.Errorf("query users error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(users, count, req.Page, req.PageSize)

	return
}

func (r *ProjectRepo) RemoveMember(userId, projectId int) (err error) {
	/*
		err = r.DB.Model(&model.ProjectMember{}).
			Where("user_id = ? AND project_id = ?", userId, projectId).
			Updates(map[string]interface{}{"deleted": true}).Error
		if err != nil {
			return
		}
	*/
	err = r.DB.
		Where("user_id = ? AND project_id=?", userId, projectId).
		Delete(&model.ProjectMember{}).Error

	return
}

func (r *ProjectRepo) FindRolesByUser(userId uint) (ret []model.ProjectMember, err error) {
	var members []model.ProjectMember

	r.DB.Model(&model.ProjectMember{}).
		Joins("LEFT JOIN biz_project_role r ON biz_project_member.project_role_id=r.id").
		Select("biz_project_member.*, r.name project_role_name").
		Where("biz_project_member.user_id = ?", userId).
		Find(&members)

	return
}

func (r *ProjectRepo) GetProjectsAndRolesByUser(userId uint) (projectIds, roleIds []uint) {
	var members []model.ProjectMember
	r.DB.Model(&model.ProjectMember{}).
		Where("user_id = ?", userId).
		Find(&members)

	roleIdsMap := make(map[uint]uint)

	for _, member := range members {
		projectIds = append(projectIds, member.ProjectId)
		roleIdsMap[member.ProjectRoleId] = member.ProjectRoleId
	}
	for _, v := range roleIdsMap {
		roleIds = append(roleIds, v)
	}

	return
}

func (r *ProjectRepo) AddProjectDefaultServe(projectId, userId uint) (serve model.Serve, err error) {
	serve = model.Serve{
		Name:      "默认服务",
		ProjectId: projectId,
	}

	err = r.DB.Create(&serve).Error

	r.ServeRepo.SetCurrServeByUser(serve.ID, userId)

	r.ServeRepo.AddDefaultServer(serve.ProjectId, serve.ID)

	return
}

func (r *ProjectRepo) FindRolesByProjectAndUser(projectId, userId uint) (projectMember model.ProjectMember, err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", projectId).
		Where("user_id = ?", userId).
		Scan(&projectMember).Error
	return
}

func (r *ProjectRepo) UpdateUserRole(req v1.UpdateProjectMemberReq) (err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", req.ProjectId).
		Where("user_id = ?", req.UserId).
		Updates(map[string]interface{}{"project_role_id": req.ProjectRoleId}).Error

	if err != nil {
		logUtils.Errorf("update project user role error", err.Error())
		return err
	}

	return
}

func (r *ProjectRepo) GetCurrProjectMemberRoleByUser(userId uint) (ret model.ProjectMember, err error) {
	curProject, err := r.GetCurrProjectByUser(userId)
	if err != nil {
		return
	}
	if curProject.ID == 0 {
		return ret, errors.New("current project is not existed")
	}
	return r.FindRolesByProjectAndUser(curProject.ID, userId)
}

func (r *ProjectRepo) GetMembersByProject(projectId uint) (ret []model.ProjectMember, err error) {
	err = r.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", projectId).
		Find(&ret).Error
	return
}

func (r *ProjectRepo) GetAuditList(req v1.AuditProjectPaginate) (data _domain.PageData, err error) {
	req.Field = "status asc,created_at"

	var count int64
	db := r.DB.Model(&model.ProjectMemberAudit{})
	if req.Type == 0 {
		projectIds := r.GetProjectIdsByUserIdAndRole(req.AuditUserId, 1)
		db = db.Where("project_id in ? and status = 0", projectIds)
	} else {
		db = db.Where("apply_user_id = ?", req.ApplyUserId)
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count ProjectMemberAudit error", zap.String("error:", err.Error()))
		return
	}

	list := make([]*model.ProjectMemberAudit, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&list).Error
	if err != nil {
		logUtils.Errorf("query ProjectMemberAudit error", zap.String("error:", err.Error()))
		return
	}

	r.refUserName(list)
	r.refProjectName(list)

	data.Populate(list, count, req.Page, req.PageSize)

	return
}

func (r *ProjectRepo) refUserName(list []*model.ProjectMemberAudit) {
	names := make(map[uint]string)
	for key, item := range list {
		if _, ok := names[item.ApplyUserId]; !ok {
			user, _ := r.UserRepo.FindById(item.ApplyUserId)
			names[item.ApplyUserId] = user.Name

		}
		if _, ok := names[item.AuditUserId]; !ok {
			user, _ := r.UserRepo.FindById(item.AuditUserId)
			names[item.AuditUserId] = user.Name
		}

		list[key].AuditUserName = names[item.AuditUserId]
		list[key].ApplyUserName = names[item.ApplyUserId]

	}
}

func (r *ProjectRepo) refProjectName(list []*model.ProjectMemberAudit) {
	names := make(map[uint]string)
	for key, item := range list {
		if _, ok := names[item.ProjectId]; !ok {
			project, _ := r.Get(item.ProjectId)
			names[item.ProjectId] = project.Name
		}

		list[key].ProjectName = names[item.ProjectId]
	}
}

func (r *ProjectRepo) GetAudit(id uint) (ret model.ProjectMemberAudit, err error) {
	err = r.DB.Model(&model.ProjectMemberAudit{}).
		Where("id = ?", id).
		First(&ret).Error
	return
}

func (r *ProjectRepo) UpdateAuditStatus(id, auditUserId uint, status consts.AuditStatus) (err error) {
	err = r.DB.Model(&model.ProjectMemberAudit{}).
		Where("id=?", id).
		Updates(map[string]interface{}{"status": status, "audit_user_id": auditUserId}).Error
	return
}

func (r *ProjectRepo) SaveAudit(audit model.ProjectMemberAudit) (err error) {
	err = r.DB.Save(&audit).Error
	return
}

func (r *ProjectRepo) IfProjectMember(userId, projectId uint) (res bool, err error) {
	var count int64
	err = r.DB.Model(&model.ProjectMember{}).Where("user_id=? and project_id=?", userId, projectId).Count(&count).Error
	if err != nil {
		return
	}
	res = count > 0
	return
}

func (r *ProjectRepo) CreateSample(projectId, serveId, userId uint) (err error) {

	//获取接口配置
	var endpoint model.Endpoint
	endpointJson := _fileUtils.ReadFile("./config/sample/endpoint.json")
	_commUtils.JsonDecode(endpointJson, &endpoint)

	//fmt.Println(_fileUtils.GetWorkDir(), endpoint)

	user, _ := r.UserRepo.FindById(userId)
	//return err

	return r.DB.Transaction(func(tx *gorm.DB) error {
		//创建接口
		endpoint.ServeId = serveId
		endpoint.ProjectId = projectId
		endpoint.CreateUser = user.Username
		err = r.EndpointRepo.saveEndpoint(&endpoint)
		if err != nil {
			return err
		}

		//TODO 创建场景
		//TODO 创建计划

		return nil
	})

}

func (r *ProjectRepo) GetProjectIdsByUserIdAndRole(userId, roleId uint) (projectIds []uint) {
	var projects []model.ProjectMember
	err := r.DB.Where("user_id=? and project_role_id=? and not deleted and not disabled", userId, roleId).Find(&projects).Error
	if err != nil {
		return
	}
	for _, project := range projects {
		projectIds = append(projectIds, project.ProjectId)
	}
	return
}
