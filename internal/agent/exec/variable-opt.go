package agentExec

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
)

// GetVariable
// @Param	processorId		int			true	"Processor Id"
// @Param	nameWithProp	string		true	"Variable name(email) or name with prop(address.city)"
// @Param	session 		ExecSession	true	"Execution Session"
func GetVariable(nameWithProp string, processorId uint, session *ExecSession) (ret domain.ExecVariable, err error) {
	/** for interface debug:
	  	    Dynamic Vars: generated by current interface's earlier pre/post conditions (agent side)
	  	    Shared Vars:  generated by other interface in same serve                   (server side)
	  								   other processor(interface) in same scenario     (server side)

	    for scenario exec:
	  	    Dynamic Vars: generated by earlier processors including interface pre/post conditions (agent side)
	  	    Shared Vars:  generated by other interface in same serve                              (server side)
	  **/

	// priority 1: Dynamic Vars, generated by execution (extractor, condition script, processor)
	ret, err = getDynamicVariableFromScope(nameWithProp, processorId, session)
	if ret.Name != "" {
		return
	}

	// priority 2: Shared Vars, for interface debug from server side
	//             A. shared vars generated by Interface Debug in same serve
	//			   B. shared vars generated by Extractor and Processor in scenario
	ret, err = getVariableFromShareVar(nameWithProp, session)
	if ret.Name != "" {
		return
	}

	// priority 3: Environment Vars, in project's serve settings
	ret, err = getVariableFromEnvVar(nameWithProp, session)
	if ret.Name != "" {
		return
	}

	// priority 4: Global Vars, on project level
	ret, err = getVariableFromGlobalVar(nameWithProp, session)
	if ret.Name != "" {
		return
	}

	return
}

// SetVariable
// @Param	processorId	int			true	"Processor Id, 0 if NOT executed by a scenario processor"
func SetVariable(processorId uint, variableName string, variableValue interface{},
	variableType consts.ExtractorResultType, scope consts.ExtractorScope, session *ExecSession) (
	newVariable domain.ExecVariable, err error) {

	scopeHierarchy := session.ScenarioDebug.ScopeHierarchy
	scopedVariables := session.ScenarioDebug.ScopedVariables

	found := false

	newVariable = domain.ExecVariable{
		Name:      variableName,
		Value:     variableValue,
		ValueType: variableType,
		Scope:     scope,
	}

	allValidIds := &[]uint{uint(0)}
	if processorId > 0 {
		allValidIds = scopeHierarchy[processorId]
	}

	if allValidIds != nil {
		for _, id := range *allValidIds {
			for i := 0; i < len(scopedVariables[id]); i++ {
				if scopedVariables[id][i].Name == variableName {
					scopedVariables[id][i] = newVariable

					found = true
					break
				}
			}
		}
	}

	if !found {
		scopedVariables[processorId] = append(scopedVariables[processorId], newVariable)
	}

	return
}

func ClearVariable(scopeId uint, variableName string, session *ExecSession) (err error) {
	scopeHierarchy := session.ScenarioDebug.ScopeHierarchy
	scopedVariables := session.ScenarioDebug.ScopedVariables

	deleteIndex := -1

	targetScopeId := uint(0)

	allValidIds := scopeHierarchy[scopeId]
	if allValidIds != nil {
		if scopeHierarchy[scopeId] != nil {
			for _, id := range *scopeHierarchy[scopeId] {
				for index, item := range scopedVariables[id] {
					if item.Name == variableName {
						deleteIndex = index
						targetScopeId = id
						break
					}
				}
			}
		}
	}

	if deleteIndex > -1 {
		if len(scopedVariables[targetScopeId]) == deleteIndex+1 {
			scopedVariables[targetScopeId] = make([]domain.ExecVariable, 0)
		} else {
			scopedVariables[targetScopeId] = append(
				scopedVariables[targetScopeId][:deleteIndex], scopedVariables[targetScopeId][(deleteIndex+1):]...)
		}
	}

	return
}
