package service

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
)

type EndpointService struct {
	EndpointRepo *repo.EndpointRepo `inject:""`
	ServeRepo    *repo.ServeRepo    `inject:""`
}

func NewEndpointService() *EndpointService {
	return &EndpointService{}
}

func (s *EndpointService) Paginate(req v1.EndpointReqPaginate) (ret _domain.PageData, err error) {
	ret, err = s.EndpointRepo.Paginate(req)
	return
}

func (s *EndpointService) Save(endpoint model.Endpoint) (res uint, err error) {
	//fmt.Println(_commUtils.JsonEncode(endpoint), "++++++", _commUtils.JsonEncode(req))
	err = s.EndpointRepo.SaveAll(&endpoint)
	return endpoint.ID, err
}

func (s *EndpointService) GetById(id uint) (res model.Endpoint) {
	res, _ = s.EndpointRepo.GetAll(id)
	return
}

func (s *EndpointService) DeleteById(id uint) (err error) {
	err = s.EndpointRepo.DeleteById(id)
	return
}

func (s *EndpointService) DisableById(id uint) (err error) {
	err = s.EndpointRepo.DisableById(id)
	return
}

func (s *EndpointService) Copy(id uint) (res uint, err error) {
	endpoint, _ := s.EndpointRepo.GetAll(id)
	s.removeIds(&endpoint)
	err = s.EndpointRepo.SaveAll(&endpoint)
	return endpoint.ID, err
}

func (s *EndpointService) removeIds(endpoint *model.Endpoint) {
	endpoint.ID = 0
	endpoint.CreatedAt = nil
	endpoint.UpdatedAt = nil
	for key, _ := range endpoint.PathParams {
		endpoint.PathParams[key].ID = 0
	}
	for key, _ := range endpoint.Interfaces {
		endpoint.Interfaces[key].ID = 0
		endpoint.Interfaces[key].RequestBody.ID = 0
		endpoint.Interfaces[key].RequestBody.SchemaItem.ID = 0
		for key1, _ := range endpoint.Interfaces[key].Headers {
			endpoint.Interfaces[key].Headers[key1].ID = 0
		}
		for key1, _ := range endpoint.Interfaces[key].Cookies {
			endpoint.Interfaces[key].Cookies[key1].ID = 0
		}
		for key1, _ := range endpoint.Interfaces[key].Params {
			endpoint.Interfaces[key].Params[key1].ID = 0
		}
		for key1, _ := range endpoint.Interfaces[key].ResponseBodies {
			endpoint.Interfaces[key].ResponseBodies[key1].ID = 0
			endpoint.Interfaces[key].ResponseBodies[key1].SchemaItem.ID = 0
			for key2, _ := range endpoint.Interfaces[key].ResponseBodies[key1].Headers {
				endpoint.Interfaces[key].ResponseBodies[key1].Headers[key2].ID = 0
			}
		}
	}

}

func (s *EndpointService) Yaml(endpoint model.Endpoint) (res interface{}) {
	serve, err := s.ServeRepo.Get(endpoint.ServeId)
	if err != nil {
		return
	}
	serve2conv := openapi.NewServe2conv(serve, []model.Endpoint{endpoint})
	res = serve2conv.ToV3()
	fmt.Println(res, "+++++++++++++++++++")
	return
}
