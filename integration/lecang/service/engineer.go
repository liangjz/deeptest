package service

import (
	"errors"
	integrationDomain "github.com/aaronchen2k/deeptest/integration/domain"
)

type EngineerService struct {
}

func (s *EngineerService) GetEngineeringOptions(baseUrl string) (ret []integrationDomain.EngineeringItem, err error) {
	token, err := new(User).GetToken(baseUrl)
	if err != nil {
		err = errors.New("您输入的环境URL地址有误")
		return
	}

	ret = new(RemoteService).LcContainerQueryAgent(token, baseUrl)

	return
}

func (s *EngineerService) GetServiceOptions(engineering, baseUrl string) (ret []integrationDomain.ServiceItem, err error) {
	token, err := new(User).GetToken(baseUrl)
	if err != nil {
		err = errors.New("您输入的环境URL地址有误")
		return
	}

	if engineering == "" {
		ret = new(RemoteService).LcAllServiceList(token, baseUrl)
	} else {
		ret = new(RemoteService).LcMlServiceQueryAgent(engineering, token, baseUrl)
	}

	return
}

func (s *EngineerService) GetAllServiceList(baseUrl string) (ret []integrationDomain.ServiceItem, err error) {
	token, err := new(User).GetToken(baseUrl)
	if err != nil {
		err = errors.New("您输入的环境URL地址有误")
		return
	}

	ret = new(RemoteService).LcAllServiceList(token, baseUrl)

	return
}
