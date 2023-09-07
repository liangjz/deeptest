package agentExec

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	scriptHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/script"
	fileUtils "github.com/aaronchen2k/deeptest/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"path/filepath"
)

var (
	MyVm      JsVm
	MyRequire *require.RequireModule

	VariableSettings []domain.ExecVariable
)

type JsVm struct {
	JsRuntime *goja.Runtime
}

func ExecScript(scriptObj *domain.ScriptBase, request domain.BaseRequest, response domain.DebugResponse) (err error) {
	VariableSettings = []domain.ExecVariable{}
	if MyVm.JsRuntime == nil {
		InitJsRuntime()
	}

	if scriptObj.Content == "" {
		return
	}

	result, err := MyVm.JsRuntime.RunString(scriptObj.Content)
	if err != nil {
		scriptObj.ResultStatus = consts.Fail
		scriptObj.Output = fmt.Sprintf("%v, ERROR: %s", result, err.Error())
		logUtils.Error(scriptObj.Output)
	} else {
		scriptObj.ResultStatus = consts.Pass
		scriptObj.Output = fmt.Sprintf("%v", result)
	}

	return
}

func InitJsRuntime() {
	registry := new(require.Registry) // registry 能夠被多个goja.Runtime共用

	MyVm.JsRuntime = goja.New()
	MyVm.JsRuntime.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	defineJsFuncs()
	defineGoFuncs()

	// load global script
	MyRequire = registry.Enable(MyVm.JsRuntime)
	pth := filepath.Join(consts.TmpDir, "deeptest.js")
	fileUtils.WriteFile(pth, scriptHelper.GetScript(scriptHelper.ScriptDeepTest))
	dt, err := MyRequire.Require(pth)
	if err != nil {
		logUtils.Info(err.Error())
		return
	}

	MyVm.JsRuntime.Set("dt", dt)
}

func defineJsFuncs() {
	MyVm.JsRuntime.Set("getDatapoolVariable", func(dpName, field, seq string) (ret interface{}) {
		rowIndex := getDatapoolRow(dpName, seq, ExecScene.Datapools)

		if ExecScene.Datapools[dpName] == nil {
			ret = "DATAPOOL_NOT_FOUND: " + dpName
			return
		}

		if rowIndex > len(ExecScene.Datapools[dpName])-1 {
			ret = "DATAPOOL_INDEX_OUT_OF_RANGE"
			return
		}

		ret = ExecScene.Datapools[dpName][rowIndex][field]
		if ret == nil {
			ret = "DATAPOOL_VARIABLE_NOT_FOUND: " + field
		}

		return
	})

	MyVm.JsRuntime.Set("getVariable", func(name string) interface{} {
		var scopeId uint
		if CurrScenarioProcessor != nil {
			scopeId = CurrScenarioProcessor.ParentId
		}
		vari, _ := GetVariable(scopeId, name)
		return vari.Value
	})
	MyVm.JsRuntime.Set("setVariable", func(name, val string) {
		var scopeId uint
		if CurrScenarioProcessor != nil {
			scopeId = CurrScenarioProcessor.ParentId
		}
		ret, err := SetVariable(scopeId, name, val, consts.Public)

		if err == nil {
			VariableSettings = append(VariableSettings, ret)
		}

		return
	})
	MyVm.JsRuntime.Set("clearVariable", func(name string) {
		var scopeId uint
		if CurrScenarioProcessor != nil {
			scopeId = CurrScenarioProcessor.ParentId
		}
		ClearVariable(scopeId, name)
	})
}

var (
	_setValueFunc func(name string, value interface{})
)

func defineGoFuncs() {
	// set data
	script := `function _setData(name, val) {
					dt[name] = val
				}`
	_, err := MyVm.JsRuntime.RunString(script)
	if err != nil {
		logUtils.Infof(err.Error())
	}

	err = MyVm.JsRuntime.ExportTo(MyVm.JsRuntime.Get("_setData"), &_setValueFunc)
}

func SetValue(name string, value interface{}) {
	_setValueFunc(name, value)
}
