package repo

import (
	"errors"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/casbin"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	model2 "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"golang.org/x/crypto/bcrypt"
	"strconv"

	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB              *gorm.DB         `inject:""`
	ProfileRepo     *ProfileRepo     `inject:""`
	RoleRepo        *RoleRepo        `inject:""`
	ProjectRepo     *ProjectRepo     `inject:""`
	ProjectRoleRepo *ProjectRoleRepo `inject:""`
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Paginate(req v1.UserReqPaginate) (data _domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model2.SysUser{})
	if len(req.Name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("获取用户总数错误", zap.String("错误:", err.Error()))
		return
	}

	users := make([]*v1.UserResp, 0)
	err = db.Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&users).Error
	if err != nil {
		logUtils.Errorf("获取用户分页数据错误", zap.String("错误:", err.Error()))
		return
	}

	// 查询用户角色
	r.GetRoles(users...)

	data.Result = users
	data.Populate(users, count, req.Page, req.PageSize)

	return
}

// getRoles
func (r *UserRepo) GetRoles(users ...*v1.UserResp) {
	var roleIds []string
	userRoleIds := make(map[uint][]string, 10)
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		user.ToString()
		userRoleId := casbin.GetRolesForUser(user.Id)
		userRoleIds[user.Id] = userRoleId
		roleIds = append(roleIds, userRoleId...)
	}

	roles, err := r.RoleRepo.FindInId(roleIds)
	if err != nil {
		logUtils.Errorf("get role get err ", zap.String("错误:", err.Error()))
	}

	for _, user := range users {
		for _, role := range roles {
			sRoleId := strconv.FormatInt(int64(role.Id), 10)
			if arr.InArrayS(userRoleIds[user.Id], sRoleId) {
				user.Roles = append(user.Roles, role.Name)
			}
		}
	}
}

func (r *UserRepo) FindByUserName(username string, ids ...uint) (v1.UserResp, error) {
	user := v1.UserResp{}
	db := r.DB.Model(&model2.SysUser{}).Where("username = ?", username)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&user).Error
	if err != nil {
		logUtils.Errorf("根据用户名查询用户错误", zap.String("用户名:", username), zap.Uints("ids:", ids), zap.String("错误:", err.Error()))
		return user, err
	}
	r.GetRoles(&user)
	return user, nil
}

func (r *UserRepo) FindPasswordByUserName(username string, ids ...uint) (v1.LoginResp, error) {
	user := v1.LoginResp{}
	db := r.DB.Model(&model2.SysUser{}).Select("id,password").Where("username = ?", username)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&user).Error
	if err != nil {
		logUtils.Errorf("根据用户名查询用户错误", zap.String("用户名:", username), zap.Uints("ids:", ids), zap.String("错误:", err.Error()))
		return user, err
	}
	return user, nil
}

func (r *UserRepo) Create(req v1.UserReq) (uint, error) {
	if _, err := r.FindByUserName(req.Username); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("用户名 %s 已经被使用", req.Username)
	}
	user := model2.SysUser{UserBase: req.UserBase, RoleIds: req.RoleIds}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logUtils.Errorf("密码加密错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	logUtils.Infof("添加用户", zap.String("hash:", req.Password), zap.ByteString("hash:", hash))

	user.Password = string(hash)
	err = r.DB.Model(&model2.SysUser{}).Create(&user).Error
	if err != nil {
		logUtils.Errorf("添加用户错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	project, err := r.AddProjectForUser(&user)
	if err != nil {
		logUtils.Errorf("添加用户项目错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	if err := r.AddProfileForUser(&user, project); err != nil {
		logUtils.Errorf("添加用户信息错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	if err := r.AddRoleForUser(&user); err != nil {
		logUtils.Errorf("添加用户角色错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	return user.ID, nil
}

func (r *UserRepo) Update(id uint, req v1.UserReq) error {
	if b, err := r.IsAdminUser(id); err != nil {
		return err
	} else if b {
		return errors.New("不能编辑超级管理员")
	}
	if _, err := r.FindByUserName(req.Username, id); !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	user := model2.SysUser{UserBase: req.UserBase}
	err := r.DB.Model(&model2.SysUser{}).Where("id = ?", id).Updates(&user).Error
	if err != nil {
		logUtils.Errorf("更新用户错误", zap.String("错误:", err.Error()))
		return err
	}

	if err := r.AddRoleForUser(&user); err != nil {
		logUtils.Errorf("添加用户角色错误", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func (r *UserRepo) IsAdminUser(id uint) (bool, error) {
	user, err := r.FindById(id)
	if err != nil {
		return false, err
	}
	return arr.InArrayS(user.Roles, serverConsts.AdminRoleName), nil
}

func (r *UserRepo) FindById(id uint) (v1.UserResp, error) {
	user := v1.UserResp{}
	err := r.DB.Model(&model2.SysUser{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		logUtils.Errorf("find user err ", zap.String("错误:", err.Error()))
		return user, err
	}

	r.GetRoles(&user)

	return user, nil
}

func (r *UserRepo) DeleteById(id uint) error {
	err := r.DB.Unscoped().Delete(&model2.SysUser{}, id).Error
	if err != nil {
		logUtils.Errorf("delete user by id get  err ", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func (r *UserRepo) AddProfileForUser(user *model2.SysUser, project model2.Project) (err error) {
	_, err = r.ProfileRepo.FindByUserId(user.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("用户 %s 信息已经被使用", user.Name)
	}

	profile := model2.SysUserProfile{UserId: user.ID, Email: "chenqi@deeptest.com", CurrProjectId: project.ID}
	err = r.DB.Create(&profile).Error
	if err != nil {
		logUtils.Errorf("添加用户错误", zap.String("错误:", err.Error()))
		return err
	}

	return
}

// AddRoleForUser add roles for user
func (r *UserRepo) AddRoleForUser(user *model2.SysUser) error {
	userId := strconv.FormatUint(uint64(user.ID), 10)
	oldRoleIds, err := casbin.Instance().GetRolesForUser(userId)
	if err != nil {
		logUtils.Errorf("获取用户角色错误", zap.String("错误:", err.Error()))
		return err
	}

	if len(oldRoleIds) > 0 {
		if _, err := casbin.Instance().DeleteRolesForUser(userId); err != nil {
			logUtils.Errorf("添加角色到用户错误", zap.String("错误:", err.Error()))
			return err
		}
	}
	if len(user.RoleIds) == 0 {
		return nil
	}

	var roleIds []string
	for _, userRoleId := range user.RoleIds {
		roleIds = append(roleIds, strconv.FormatUint(uint64(userRoleId), 10))
	}

	if _, err := casbin.Instance().AddRolesForUser(userId, roleIds); err != nil {
		logUtils.Errorf("添加角色到用户错误", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func (r *UserRepo) AddProjectForUser(user *model2.SysUser) (project model2.Project, err error) {
	_, err = r.ProjectRepo.GetCurrProjectByUser(user.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		err = fmt.Errorf("用户%s的默认项目已存在", user.Name)
		return
	}

	project = model2.Project{ProjectBase: v1.ProjectBase{Name: "默认项目"}}
	err = r.DB.Create(&project).Error
	if err != nil {
		logUtils.Errorf("添加项目错误", zap.String("错误:", err.Error()))
		return
	}

	err = r.ProjectRepo.AddProjectMember(project.ID, user.ID)
	if err != nil {
		logUtils.Errorf("添加项目角色错误", zap.String("错误:", err.Error()))
		return
	}

	err = r.ProjectRepo.AddProjectRootInterface(project.ID)
	if err != nil {
		logUtils.Errorf("添加接口错误", zap.String("错误:", err.Error()))
		return
	}

	return
}

// DelToken 删除token
func (r *UserRepo) DelToken(token string) error {
	err := multi.AuthDriver.DelUserTokenCache(token)
	if err != nil {
		logUtils.Errorf("del token", zap.Any("err", err))
		return fmt.Errorf("del token %w", err)
	}
	return nil
}

// CleanToken 清空 token
func (r *UserRepo) CleanToken(authorityType int, userId string) error {
	err := multi.AuthDriver.CleanUserTokenCache(authorityType, userId)
	if err != nil {
		logUtils.Errorf("clean token", zap.Any("err", err))
		return fmt.Errorf("clean token %w", err)
	}
	return nil
}

func (r *UserRepo) UpdatePasswordByName(name string, password string) (err error) {
	err = r.DB.Model(&model2.SysUser{}).Where("username = ?", name).
		Updates(map[string]interface{}{"password": password}).Error
	if err != nil {
		logUtils.Errorf("更新用户错误", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}
func (r *UserRepo) UpdateAvatar(id uint, avatar string) error {
	return nil
}