package business

import (
	"errors"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"strings"
	"time"
)

var (
	ExecLog = domain.ExecLog{}

	ScopeHierarchy  = map[uint]*[]uint{}
	ScopedVariables = map[uint][]domain.ExecVariable{}
	ScopedCookies   = map[uint][]domain.ExecCookie{}
)

type ExecContext struct {
	ScenarioNodeRepo *repo.ScenarioNodeRepo `inject:""`
}

func (s *ExecContext) InitScopeHierarchy(scenarioId uint) (variables []domain.ExecVariable) {
	s.ScenarioNodeRepo.GetScopeHierarchy(scenarioId, &ScopeHierarchy)

	ScopedVariables = map[uint][]domain.ExecVariable{}
	ScopedCookies = map[uint][]domain.ExecCookie{}

	return
}

func (s *ExecContext) ListCachedVariable(scopeId uint) (variables []domain.ExecVariable) {
	effectiveScopeIds := ScopeHierarchy[scopeId]

	for _, id := range *effectiveScopeIds {
		for _, vari := range ScopedVariables[id] {
			if !vari.IsShare && id != scopeId {
				continue
			}

			variables = append(variables, vari)
		}
	}

	return
}

func (s *ExecContext) GetVariable(scopeId uint, variablePath string) (variable domain.ExecVariable, err error) {
	if variablePath == "var1" {
		logUtils.Info("")
	}

	allValidIds := ScopeHierarchy[scopeId]
	if allValidIds != nil {
		for _, id := range *allValidIds {
			for _, item := range ScopedVariables[id] {
				var ok bool
				if variable, ok = s.EvaluateVariableExpressionValue(item, variablePath); ok {
					goto LABEL
				}
			}
		}
	}

	if variable.Name == "" { // not found
		err = errors.New(fmt.Sprintf("找不到变量\"%s\"", variablePath))
	}

LABEL:

	return
}

func (s *ExecContext) EvaluateVariableExpressionValue(variable domain.ExecVariable, variablePath string) (
	ret domain.ExecVariable, ok bool) {
	arr := strings.Split(variablePath, ".")
	variableName := arr[0]

	if variable.Name == variableName {
		ret = variable

		if len(arr) > 1 {
			variableProp := arr[1]
			ret.Value = variable.Value.(map[string]interface{})[variableProp]
		}

		ok = true

	}

	return
}

func (s *ExecContext) SetVariable(scopeId uint, variableName string, variableValue interface{}, isShare bool) (
	err error) {

	found := false

	newVariable := domain.ExecVariable{
		Name:    variableName,
		Value:   variableValue,
		IsShare: isShare,
	}

	allValidIds := ScopeHierarchy[scopeId]
	if allValidIds != nil {
		for _, id := range *allValidIds {
			for i := 0; i < len(ScopedVariables[id]); i++ {
				if ScopedVariables[id][i].Name == variableName {
					ScopedVariables[id][i] = newVariable

					found = true
					break
				}
			}
		}
	}

	if !found {
		ScopedVariables[scopeId] = append(ScopedVariables[scopeId], newVariable)
	}

	return
}

func (s *ExecContext) ClearVariable(scopeId uint, variableName string) (err error) {
	deleteIndex := -1

	targetScopeId := uint(0)

	allValidIds := ScopeHierarchy[scopeId]
	if allValidIds != nil {
		for _, id := range *ScopeHierarchy[scopeId] {
			for index, item := range ScopedVariables[id] {
				if item.Name == variableName {
					deleteIndex = index
					targetScopeId = id
					break
				}
			}
		}
	}

	if deleteIndex > -1 {
		ScopedVariables[scopeId] = append(
			ScopedVariables[targetScopeId][:deleteIndex], ScopedVariables[scopeId][(deleteIndex+1):]...)
	}

	return
}

func (s *ExecContext) ListCookie(scopeId uint) (cookies []domain.ExecCookie) {
	allValidIds := ScopeHierarchy[scopeId]
	if allValidIds != nil {
		for _, id := range *ScopeHierarchy[scopeId] {
			cookies = append(cookies, ScopedCookies[id]...)
		}
	}

	return
}

func (s *ExecContext) GetCookie(scopeId uint, cookieName, domain string) (cookie domain.ExecCookie) {
	allValidIds := ScopeHierarchy[scopeId]
	if allValidIds != nil {
		for _, id := range *ScopeHierarchy[scopeId] {
			for _, item := range ScopedCookies[id] {
				if item.Name == cookieName && item.Domain == domain && item.ExpireTime.Unix() > time.Now().Unix() {
					cookie = item

					goto LABEL
				}
			}
		}
	}

LABEL:

	return
}

func (s *ExecContext) SetCookie(scopeId uint, cookieName string, cookieValue interface{}, domainName string, expireTime *time.Time) (err error) {
	found := false

	newCookie := domain.ExecCookie{
		Name:  cookieName,
		Value: cookieValue,

		Domain:     domainName,
		ExpireTime: expireTime,
	}

	for i := 0; i < len(ScopedCookies[scopeId]); i++ {
		if ScopedCookies[scopeId][i].Name == cookieName {
			ScopedCookies[scopeId][i] = newCookie

			found = true
			break
		}
	}

	if !found {
		ScopedCookies[scopeId] = append(ScopedCookies[scopeId], newCookie)
	}

	return
}

func (s *ExecContext) ClearCookie(scopeId uint, cookieName string) (err error) {
	deleteIndex := -1
	for index, item := range ScopedCookies[scopeId] {
		if item.Name == cookieName {
			deleteIndex = index
			break
		}
	}

	if deleteIndex > -1 {
		ScopedCookies[scopeId] = append(
			ScopedCookies[scopeId][:deleteIndex], ScopedCookies[scopeId][(deleteIndex+1):]...)
	}

	return
}