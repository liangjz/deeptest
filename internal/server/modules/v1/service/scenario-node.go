package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
)

type ScenarioNodeService struct {
	ScenarioNodeRepo      *repo.ScenarioNodeRepo      `inject:""`
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
	ScenarioRepo          *repo.ScenarioRepo          `inject:""`
}

func (s *ScenarioNodeService) GetTree(scenarioId int) (root *model.Processor, err error) {
	root, err = s.ScenarioNodeRepo.GetTree(scenarioId)

	return
}

func (s *ScenarioNodeService) AddInterfaces(req serverDomain.ScenarioAddInterfacesReq) (ret model.Processor, err *_domain.BizErr) {
	targetProcessor, _ := s.ScenarioProcessorRepo.Get(req.TargetId)

	for _, interfaceNode := range req.SelectedNodes {
		ret, _ = s.createDirOrInterface(interfaceNode, targetProcessor)
	}

	return
}

func (s *ScenarioNodeService) AddProcessor(req serverDomain.ScenarioAddScenarioReq) (ret model.Processor, err *_domain.BizErr) {
	targetProcessor, _ := s.ScenarioProcessorRepo.Get(uint(req.TargetProcessorId))
	if targetProcessor.ID == 0 {
		return
	}

	ret = model.Processor{
		Name:           req.Name,
		EntityCategory: req.ProcessorCategory,
		EntityType:     req.ProcessorType,

		ScenarioId: targetProcessor.ScenarioId,
		ProjectId:  req.ProjectId,
	}

	if req.Mode == "child" {
		ret.ParentId = targetProcessor.ID
	} else if req.Mode == "parent" && req.TargetProcessorCategory == consts.ProcessorInterface {
		ret.ParentId = targetProcessor.ParentId
	}
	ret.Ordr = s.ScenarioNodeRepo.GetMaxOrder(ret.ParentId)

	s.ScenarioNodeRepo.Save(&ret)

	if req.Mode == "parent" {
		targetProcessor.ParentId = ret.ID
		s.ScenarioNodeRepo.Save(&targetProcessor)
	}

	if ret.EntityType == consts.ProcessorLogicElse { // create default entity
		entity := model.ProcessorLogic{
			ProcessorEntity: model.ProcessorEntity{
				ProcessorId:       ret.ID,
				ProcessorCategory: ret.EntityCategory,
				ProcessorType:     ret.EntityType,
			},
		}
		s.ScenarioProcessorRepo.SaveLogic(&entity)
	}

	return
}

func (s *ScenarioNodeService) createDirOrInterface(interfaceNode serverDomain.InterfaceSimple, parentProcessor model.Processor) (
	ret model.Processor, err *_domain.BizErr) {

	if interfaceNode.IsDir {
		processor := model.Processor{
			Name:           interfaceNode.Name,
			ScenarioId:     parentProcessor.ScenarioId,
			EntityCategory: consts.ProcessorGroup,
			ParentId:       parentProcessor.ID,
			ProjectId:      parentProcessor.ProjectId,
		}
		processor.Ordr = s.ScenarioNodeRepo.GetMaxOrder(processor.ParentId)
		s.ScenarioNodeRepo.Save(&processor)

		for _, child := range interfaceNode.Children {
			s.createDirOrInterface(child, processor)
		}

	} else {
		processor := model.Processor{
			Name:           interfaceNode.Name,
			ScenarioId:     parentProcessor.ScenarioId,
			EntityCategory: consts.ProcessorInterface,
			EntityType:     consts.ProcessorInterfaceDefault,
			InterfaceId:    uint(interfaceNode.Id),
			ParentId:       parentProcessor.ID,
			ProjectId:      parentProcessor.ProjectId,
		}
		processor.Ordr = s.ScenarioNodeRepo.GetMaxOrder(processor.ParentId)
		s.ScenarioNodeRepo.Save(&processor)

		ret = processor
	}

	return
}

func (s *ScenarioNodeService) UpdateName(req serverDomain.ScenarioNodeReq) (err error) {
	err = s.ScenarioNodeRepo.UpdateName(req.Id, req.Name)
	return
}

func (s *ScenarioNodeService) Delete(id uint) (err error) {
	err = s.deleteScenarioNodeAndChildren(id)

	return
}

func (s *ScenarioNodeService) Move(srcId, targetId uint, pos serverConsts.DropPos, projectId uint) (
	srcScenarioNode model.Processor, err error) {
	srcScenarioNode, err = s.ScenarioNodeRepo.Get(srcId)

	srcScenarioNode.ParentId, srcScenarioNode.Ordr = s.ScenarioNodeRepo.UpdateOrder(pos, targetId)
	err = s.ScenarioNodeRepo.UpdateOrdAndParent(srcScenarioNode)

	return
}

func (s *ScenarioNodeService) deleteScenarioNodeAndChildren(nodeId uint) (err error) {
	err = s.ScenarioNodeRepo.Delete(nodeId)
	if err == nil {
		children, _ := s.ScenarioNodeRepo.GetChildren(nodeId)
		for _, child := range children {
			s.deleteScenarioNodeAndChildren(child.ID)
		}
	}

	return
}
