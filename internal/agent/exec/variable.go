package agentExec

import (
	"errors"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"strings"
)

//func ImportVariables(processorId uint, variables domain.VarKeyValuePair, scope consts.ExtractorScope) (err error) {
//	for key, val := range variables {
//		newVariable := domain.ExecVariable{
//			Name:  key,
//			Value: val,
//			Scope: scope,
//		}
//
//		found := false
//		for i := 0; i < len(ScopedVariables[processorId]); i++ {
//			if ScopedVariables[processorId][i].Name == key {
//				ScopedVariables[processorId][i] = newVariable
//
//				found = true
//				break
//			}
//		}
//
//		if !found {
//			ScopedVariables[processorId] = append(ScopedVariables[processorId], newVariable)
//		}
//	}
//
//	return
//}

func GetVariable(processorId uint, variablePath string) (variable domain.ExecVariable, err error) {
	allValidIds := ScopeHierarchy[processorId]
	if allValidIds != nil {
		for _, id := range *allValidIds {
			for _, item := range ScopedVariables[id] {
				var ok bool
				if variable, ok = EvaluateVariableExpressionValue(item, variablePath); ok {
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

func SetVariable(processorId uint, variableName string, variableValue interface{}, scope consts.ExtractorScope) (
	err error) {

	found := false

	newVariable := domain.ExecVariable{
		Name:  variableName,
		Value: variableValue,
		Scope: scope,
	}

	allValidIds := ScopeHierarchy[processorId]
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
		ScopedVariables[processorId] = append(ScopedVariables[processorId], newVariable)
	}

	return
}

func ClearVariable(processorId uint, variableName string) (err error) {
	deleteIndex := -1

	targetScopeId := uint(0)

	allValidIds := ScopeHierarchy[processorId]
	if allValidIds != nil {
		for _, id := range *ScopeHierarchy[processorId] {
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
		if len(ScopedVariables[targetScopeId]) == deleteIndex+1 {
			ScopedVariables[targetScopeId] = make([]domain.ExecVariable, 0)
		} else {
			ScopedVariables[targetScopeId] = append(
				ScopedVariables[targetScopeId][:deleteIndex], ScopedVariables[targetScopeId][(deleteIndex+1):]...)
		}
	}

	return
}

func ReplaceVariableValue(value string) (ret string) {
	variablePlaceholders := GetVariablesInVariablePlaceholder(value)
	ret = value

	for _, placeholder := range variablePlaceholders {
		oldVal := fmt.Sprintf("${%s}", placeholder)

		newVal := getPlaceholderValue(placeholder)

		ret = strings.ReplaceAll(ret, oldVal, newVal)
	}

	return
}

func getPlaceholderValue(name string) (ret string) {
	typ := getPlaceholderType(name)

	if typ == consts.PlaceholderTypeVariable {
		ret = getVariableValue(name)

	} else if typ == consts.PlaceholderTypeDatapool {
		ret = getDatapoolValue(name)

	}
	//else if typ == consts.PlaceholderTypeFunction {
	//}

	return
}

func getVariableValue(name string) (ret string) {
	// priority 1: A. shared vars generated by Debug Endpoint Interface in same serve
	//			   B. shared vars generated by Extractor and Processor in scenario
	ret = getValueFromShareVar(name)
	if ret != "" {
		return
	}

	// priority 2: environment vars on project's serve settings
	ret = getValueFromEnvVar(name)
	if ret != "" {
		return
	}

	// priority 3: global vars on project level
	ret = getValueFromGlobalVar(name)
	if ret != "" {
		return
	}

	return
}
func getValueFromShareVar(name string) (ret string) {
	if CurrProcessorId == 0 { // endpoint interface dbug
		ret = getValueFromList(name, ExecScene.ShareVars)

	} else { // run scenario
		cache := CachedShareVarByProcessor[CurrProcessorId]
		if cache == nil {
			cache = GetCachedVariableMapInContext(CurrProcessorId)
		}

		if cache[name] == nil {
			return ""
		}

		ret = fmt.Sprintf("%v", cache[name])
	}

	return
}
func getValueFromEnvVar(name string) (ret string) {
	envId := ExecScene.InterfaceToEnvMap[CurrInterfaceId]

	vars := ExecScene.EnvToVariables[envId]

	ret = getValueFromList(name, vars)

	return
}
func getValueFromGlobalVar(name string) (ret string) {
	ret = getValueFromList(name, ExecScene.GlobalVars)

	return
}

func getValueFromList(name string, list []domain.GlobalVar) (ret string) {
	for _, v := range list {
		if v.Name == name {
			ret = v.LocalValue
			break
		}
	}

	return
}

func getPlaceholderType(placeholder string) (ret consts.PlaceholderType) {
	if strings.HasPrefix(placeholder, consts.PlaceholderPrefixDatapool.String()) {
		return consts.PlaceholderTypeDatapool
	} else if strings.HasPrefix(placeholder, consts.PlaceholderPrefixFunction.String()) {
		return consts.PlaceholderTypeFunction
	}

	return consts.PlaceholderTypeVariable
}