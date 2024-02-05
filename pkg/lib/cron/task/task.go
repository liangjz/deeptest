package task

import "github.com/aaronchen2k/deeptest/internal/pkg/core/cron"

type Task interface {
	Run(options map[string]interface{}) (f func() error)
	CallBack(options map[string]interface{}, err error) func()
}

type Proxy struct {
	source string
	cron   string
	task   Task
	taskId string
	Cron   *cron.ServerCron `inject:""`
}

func (p *Proxy) GetTaskId() (taskId string) {
	taskId = p.source + "_" + p.taskId
	return
}

func NewProxy(source, cron string) (proxy Proxy) {
	proxy = Proxy{
		source: source,
		cron:   cron,
	}

	taskEntity := Factory{
		name: source,
	}

	proxy.task = taskEntity.Create()
	return
}

func (p *Proxy) Add(options map[string]interface{}) (err error) {
	taskFunc := p.getTaskFunc(options)

	err = p.Cron.AddCommonTask(p.GetTaskId(), p.cron, taskFunc)

	return
}

func (p *Proxy) getTaskFunc(options map[string]interface{}) (taskFunc func()) {
	runFunc := p.task.Run(options)
	if runFunc() == nil {
		return
	}

	taskFunc = func() {
		err := runFunc()
		callBackFunc := p.task.CallBack(options, err)
		if callBackFunc == nil {
			return
		}

		callBackFunc()
	}

	return
}

func Test() {
	options := make(map[string]interface{})
	options["swagger_1"] = 1
	proxy := NewProxy("swagger", "*****")
	proxy.Add(options)
}
