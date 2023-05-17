package repo

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"gorm.io/gorm"
	"strconv"
)

type EndpointRepo struct {
	*BaseRepo              `inject:""`
	EndpointInterfaceRepo  *EndpointInterfaceRepo  `inject:""`
	ServeRepo              *ServeRepo              `inject:""`
	ProcessorInterfaceRepo *ProcessorInterfaceRepo `inject:""`
	ProjectRepo            *ProjectRepo            `inject:""`
}

func NewEndpointRepo() *EndpointRepo {
	return &EndpointRepo{}
}

func (r *EndpointRepo) Paginate(req v1.EndpointReqPaginate) (ret _domain.PageData, err error) {
	//fmt.Println(r.DB.Model(&model.SysUser{}))
	//err = r.DB.Where("id=?", id).Where("name=?", name).Find(&res).Error
	var count int64
	db := r.DB.Model(&model.Endpoint{}).Where("project_id = ? AND NOT deleted AND NOT disabled", req.ProjectId)

	if req.Title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Title))
	}
	if req.CreateUser != "" {
		db = db.Where("create_user = ?", req.CreateUser)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.ServeId != 0 {
		db = db.Where("serve_id = ?", req.ServeId)
	}
	if req.ServeVersion != "" {
		if ids, err := r.ServeRepo.GetBindEndpointIds(req.ServeId, req.ServeVersion); err != nil {
			db = db.Where("id in ?", ids)
		}
	}

	if req.CategoryId > 0 {
		var categoryIds []uint
		categoryIds, err = r.BaseRepo.GetAllChildIds(uint(req.CategoryId), model.Category{}.TableName(),
			serverConsts.EndpointCategory, int(req.ProjectId))
		if err != nil {
			return
		}
		if len(categoryIds) > 0 {
			db.Where("category_id IN(?)", categoryIds)
		}
	} else if req.CategoryId == -1 {
		db.Where("category_id IN(-1)")
	}

	db = db.Order("created_at desc")
	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count report error %s", err.Error())
		return
	}

	results := make([]*model.Endpoint, 0)

	err = db.Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).Find(&results).Error
	if err != nil {
		logUtils.Errorf("query report error %s", err.Error())
		return
	}

	serveNames := map[uint]string{}

	for key, result := range results {
		var versions []model.EndpointVersion
		r.DB.Find(&versions, "endpoint_id=?", result.ID).Order("version desc")
		results[key].Versions = versions
		if len(versions) > 0 {
			results[key].Version = versions[0].Version
		}

		if _, ok := serveNames[result.ServeId]; !ok {
			var serve model.Serve
			r.DB.Find(&serve, "id=?", result.ServeId)
			serveNames[result.ServeId] = serve.Name
		}
		results[key].ServeName = serveNames[result.ServeId]
	}

	ret.Populate(results, count, req.Page, req.PageSize)

	return
}

func (r *EndpointRepo) SaveAll(endpoint *model.Endpoint) (err error) {
	r.DB.Transaction(func(tx *gorm.DB) error {
		//创建version
		err = r.saveEndpointVersion(endpoint)
		if err != nil {
			return err
		}

		//更新终端
		err = r.saveEndpoint(endpoint)
		if err != nil {
			return err
		}

		//保存路径参数
		err = r.saveEndpointParams(endpoint.ID, endpoint.PathParams)
		if err != nil {
			return err
		}

		//保存接口
		err = r.saveInterfaces(endpoint.ID, endpoint.ProjectId, endpoint.Path, endpoint.Version, endpoint.Interfaces)
		if err != nil {
			return err
		}

		//更新调试接口
		//err = r.updateDebugInterfaceUrl(endpoint.ID, endpoint.Path)
		//if err != nil {
		//	return err
		//}

		return nil
	})
	return
}

//保存终端信息
func (r *EndpointRepo) saveEndpoint(endpoint *model.Endpoint) (err error) {
	err = r.Save(endpoint.ID, endpoint)
	if err != nil {
		return
	}

	err = r.UpdateSerialNumber(endpoint.ID, uint(endpoint.ProjectId))
	return
}

func (r *EndpointRepo) UpdateSerialNumber(id, projectId uint) (err error) {
	var project model.Project
	project, err = r.ProjectRepo.Get(projectId)
	if err != nil {
		return
	}

	err = r.DB.Model(&model.Endpoint{}).Where("id=?", id).Update("serial_number", project.ShortName+"-I-"+strconv.Itoa(int(id))).Error
	return

}

func (r *EndpointRepo) saveEndpointVersion(endpoint *model.Endpoint) (err error) {
	if endpoint.Version == "" {
		endpoint.Version = "v0.1.0"
	}

	endpointVersion := model.EndpointVersion{EndpointId: endpoint.ID, Version: endpoint.Version}
	r.FindVersion(&endpointVersion)
	if endpointVersion.ID == 0 {
		err = r.DB.Create(&endpointVersion).Error
		if err != nil {
			endpoint.Version = endpointVersion.Version
		}
	}

	return
}

//保存路径参数
func (r *EndpointRepo) saveEndpointParams(endpointId uint, params []model.EndpointPathParam) (err error) {
	err = r.removeEndpointParams(endpointId)
	if err != nil {
		return
	}
	for _, item := range params {
		item.EndpointId = endpointId
		err = r.Save(item.ID, &item)
		if err != nil {
			return
		}
	}
	return
}

func (r *EndpointRepo) removeEndpointParams(endpointId uint) (err error) {
	err = r.DB.
		Where("endpoint_id = ?", endpointId).
		Delete(&model.EndpointPathParam{}, "").Error

	return
}

//保存接口信息
func (r *EndpointRepo) saveInterfaces(endpointId, projectId uint, path, version string, interfaces []model.EndpointInterface) (err error) {
	err = r.removeInterfaces(endpointId)
	if err != nil {
		return
	}
	for _, item := range interfaces {
		item.EndpointId = endpointId
		item.Version = version
		item.Url = path
		item.ProjectId = projectId

		err = r.EndpointInterfaceRepo.SaveInterfaces(&item)
		if err != nil {
			return err
		}
	}
	return
}

//保存调试接口Url
func (r *EndpointRepo) updateDebugInterfaceUrl(endpointId uint, url string) (err error) {
	err = r.DB.Model(&model.DebugInterface{}).
		Where("endpoint_id = ?", endpointId).
		Update("url", url).Error

	return
}

func (r *EndpointRepo) removeInterfaces(endpointId uint) (err error) {
	err = r.DB.
		Where("endpoint_id = ?", endpointId).
		Delete(&model.EndpointInterface{}, "").Error

	return
}

func (r *EndpointRepo) GetAll(id uint, version string) (endpoint model.Endpoint, err error) {
	endpoint, err = r.Get(id)
	if err != nil {
		return
	}
	endpoint.PathParams, _ = r.GetEndpointParams(id)
	endpoint.Interfaces, _ = r.EndpointInterfaceRepo.GetByEndpointId(id, version)

	return
}

func (r *EndpointRepo) Get(id uint) (res model.Endpoint, err error) {
	err = r.DB.First(&res, id).Error
	return
}

func (r *EndpointRepo) GetEndpointParams(endpointId uint) (pathParam []model.EndpointPathParam, err error) {
	err = r.DB.Find(&pathParam, "endpoint_id=?", endpointId).Error
	return
}

func (r *EndpointRepo) DeleteById(id uint) error {
	return r.DB.Model(&model.Endpoint{}).Where("id = ?", id).Update("deleted", 1).Error
}

func (r *EndpointRepo) DisableById(id uint) error {
	return r.DB.Model(&model.Endpoint{}).Where("id = ?", id).Update("status", 4).Error
}

func (r *EndpointRepo) UpdateStatus(id uint, status int64) error {
	return r.DB.Model(&model.Endpoint{}).Where("id = ?", id).Update("status", status).Error
}

func (r *EndpointRepo) DeleteByIds(ids []uint) error {
	return r.DB.Model(&model.Endpoint{}).Where("id IN ?", ids).Update("deleted", 1).Error
}

func (r *EndpointRepo) GetVersionsByEndpointId(endpointId uint) (res []model.EndpointVersion, err error) {
	err = r.DB.Find(&res, "endpoint_id=?", endpointId).Error
	return
}

func (r *EndpointRepo) GetLatestVersion(endpointId uint) (res model.EndpointVersion, err error) {
	err = r.DB.Take(&res, "endpoint_id=?", endpointId).Order("version desc").Error
	return
}
func (r *EndpointRepo) FindVersion(res *model.EndpointVersion) (err error) {
	err = r.DB.Where("endpoint_id=? and version=?", res.EndpointId, res.Version).First(&res).Error
	return
}

func (r *EndpointRepo) GetFirstMethod(id uint) (res model.EndpointInterface, err error) {
	var interfs []model.EndpointInterface
	interfs, err = r.EndpointInterfaceRepo.GetByEndpointId(id, "v0.1.0")
	if len(interfs) > 0 {
		res = interfs[0]
	}
	return
}

func (r *EndpointRepo) GetCountByServeId(serveId uint) (count int64, err error) {
	err = r.DB.Where("serve=? and NOT deleted", serveId).Count(&count).Error
	return
}

func (r *EndpointRepo) ListEndpointByCategory(categoryId uint) (ids []uint, err error) {
	err = r.DB.Model(&model.Endpoint{}).
		Select("id").
		Where("category_id = ? AND NOT deleted", categoryId).
		Find(&ids).Error
	return
}

func (r *EndpointRepo) CreateEndpointSample(serveId uint) (endpointId uint, err error) {

	return
}
