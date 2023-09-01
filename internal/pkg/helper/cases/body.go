package cases

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/getkin/kin-openapi/openapi3"
)

func GenerateByBody(basicDebugData domain.DebugData, apiOperation *openapi3.Operation) (
	alternativeCase []model.EndpointCase, err error) {

	return
}