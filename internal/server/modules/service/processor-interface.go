package service

import (
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/jinzhu/copier"
)

type ProcessorInterfaceService struct {
	ProcessorInterfaceRepo *repo.ProcessorInterfaceRepo `inject:""`
}

func NewScenarioInterfaceService() *ProcessorInterfaceService {
	return &ProcessorInterfaceService{}
}

func (s *ProcessorInterfaceService) GetById(id uint) (model.ProcessorInterface, error) {
	return s.ProcessorInterfaceRepo.GetDetail(id)
}

func (s *ProcessorInterfaceService) ListInvocation(id uint) (invocations []model.ProcessorInvocation, err error) {
	return s.ProcessorInterfaceRepo.ListInvocation(id)
}

func (s *ProcessorInterfaceService) UpdateByInvocation(req v1.DebugData) (err error) {
	interf := model.ProcessorInterface{}
	s.CopyValueFromRequest(&interf, req)

	err = s.ProcessorInterfaceRepo.Update(interf)

	return
}

func (s *ProcessorInterfaceService) CopyValueFromRequest(interf *model.ProcessorInterface, req v1.DebugData) (err error) {
	interf.ID = req.EndpointInterfaceId

	copier.CopyWithOption(interf, req, copier.Option{DeepCopy: true})

	return
}

func (s *ProcessorInterfaceService) UpdateByConfig(req v1.DebugData) (err error) {
	interf := model.ProcessorInterface{}
	s.CopyValueFromRequest(&interf, req)

	err = s.ProcessorInterfaceRepo.Update(interf)

	return
}

func (s *ProcessorInterfaceService) GetList(projectId, scenarioId uint) (res []model.ProcessorInterface, err error) {
	res, err = s.ProcessorInterfaceRepo.GetList(projectId, scenarioId)
	return
}
