package agentExec

type PlanExecReq struct {
	PlanId        uint   `json:"planId"`
	EnvId         uint   `json:"envId"`
	ServerUrl     string `json:"serverUrl"`
	Token         string `json:"token"`
	EnvironmentId int    `json:"environmentId"`
}

type PlanExecObj struct {
	Name      string            `json:"name"`
	Scenarios []ScenarioExecObj `json:"scenarios"`

	ServerUrl string `json:"serverUrl"`
	Token     string `json:"token"`
}
