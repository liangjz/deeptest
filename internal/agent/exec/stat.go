package agentExec

import (
	"encoding/json"
	agentDomain "github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"regexp"
	"strings"
)

func ResetStat(session ExecSession) {
	session.InterfaceStat = agentDomain.InterfaceStat{}
}

func CountInterfaceStat(result *agentDomain.ScenarioExecResult, execUuid string) agentDomain.InterfaceStat {
	stat := GetInterfaceStat(execUuid)

	stat.InterfaceCount += 1
	stat.InterfaceDurationTotal += result.Cost
	stat.InterfaceDurationAverage = stat.InterfaceDurationTotal / stat.InterfaceCount

	result.ResultStatus = consts.Pass

	for _, item := range result.PostConditions {
		if item.Type != consts.ConditionTypeCheckpoint {
			continue
		}

		var checkpointBase domain.CheckpointBase
		json.Unmarshal(item.Raw, &checkpointBase)
		if checkpointBase.Disabled {
			continue
		}

		if checkpointBase.ResultStatus == consts.Pass {
			stat.CheckpointPass += 1
		} else if checkpointBase.ResultStatus == consts.Fail {
			stat.CheckpointFail += 1
			result.ResultStatus = consts.Fail
		}
		/*
			if item.Type != consts.ConditionTypeResponseDefine {
				var responseDefineBase domain.ResponseDefineBase
				json.Unmarshal(item.Raw, &responseDefineBase)
				if responseDefineBase.Disabled {
					continue
				}
			}
		*/
	}

	if result.ResultStatus == consts.Pass {
		stat.InterfacePass += 1
	} else if result.ResultStatus == consts.Fail {
		stat.InterfaceFail += 1
	}

	SetInterfaceStat(execUuid, stat)

	return *stat
}

func CountScriptAssertionStat(output string, result *agentDomain.ScenarioExecResult, execUuid string) agentDomain.InterfaceStat {
	stat := GetInterfaceStat(execUuid)

	arr := []string{}
	json.Unmarshal([]byte(output), &arr)

	for _, item := range arr {
		status, _, _ := ParseChaiAssertion(item)
		if status == "pass" {
			result.Stat.CheckpointPass += 1
			stat.CheckpointPass += 1

		} else if status == "fail" {
			result.Stat.CheckpointFail += 1
			stat.CheckpointFail += 1
		}
	}

	SetInterfaceStat(execUuid, stat)

	return *stat
}

func CountSkip(executedProcessorIds map[uint]bool, skippedChildren []*Processor, session *ExecSession) agentDomain.InterfaceStat {
	countedProcessorIds := map[uint]bool{}
	countSkipInterface(executedProcessorIds, skippedChildren, &countedProcessorIds, session)

	return session.InterfaceStat
}

func countSkipInterface(executedProcessorIds map[uint]bool, skippedChildren []*Processor, countedProcessorIds *map[uint]bool, session *ExecSession) agentDomain.InterfaceStat {
	stat := session.InterfaceStat

	for _, child := range skippedChildren {
		if child.Disable {
			continue
		}

		_, executed := executedProcessorIds[child.ID]
		_, counted := (*countedProcessorIds)[child.ID]
		if child.EntityType == consts.ProcessorInterfaceDefault && !executed && !counted {
			stat.InterfaceSkip += 1
			(*countedProcessorIds)[child.ID] = true
		}

		if len(child.Children) > 0 {
			countSkipInterface(map[uint]bool{}, child.Children, countedProcessorIds, session)
		}
	}

	session.InterfaceStat = stat
	return session.InterfaceStat
}

func ParseChaiAssertion(output string) (status, name, checkpoint string) {
	// Assertion Pass [Assertion 1].
	// Assertion Failed [Assertion 1] AssertionError: check status code: expected 200 to equal 2001.

	regx := regexp.MustCompile(`Assertion (Failed|Pass) \[(.+)\](.*)\.`)
	arr := regx.FindAllStringSubmatch(output, -1)
	if len(arr) == 0 {
		return
	}

	status = strings.ToLower(arr[0][1])
	if status != "pass" {
		status = "fail"
	}
	name = arr[0][2]
	checkpoint = arr[0][3]

	return
}
