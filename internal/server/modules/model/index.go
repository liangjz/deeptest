package model

var (
	Models = []interface{}{
		&Oplog{},

		&SysPerm{},
		&SysRole{},
		&SysUser{},
		&SysUserProfile{},

		&ProjectRole{},
		&Org{},
		&Project{},
		&ProjectMember{},
		&Datapool{},
		&Environment{},
		&EnvironmentVar{},
		&ShareVariable{},

		&DebugInterface{},
		&DebugInterfaceParam{},
		&DebugInterfaceBodyFormDataItem{},
		&DebugInterfaceBodyFormUrlEncodedItem{},
		&DebugInterfaceHeader{},
		&DebugInterfaceBasicAuth{},
		&DebugInterfaceBearerToken{},
		&DebugInterfaceOAuth20{},
		&DebugInterfaceApiKey{},

		&DebugCondition{},
		&DebugConditionExtractor{},
		&DebugConditionCheckpoint{},
		&DebugConditionScript{},

		&DiagnoseInterface{},

		&Snippet{},
		&MockJsExpression{},

		&MockInvocation{},
		&Auth2Token{},

		&Category{},
		&Scenario{},

		&Plan{},
		&RelaPlanScenario{},

		&Processor{},
		//&ProcessorThreadGroup{},
		&ProcessorGroup{},
		&ProcessorLogic{},
		&ProcessorLoop{},
		&ProcessorTimer{},
		&ProcessorPrint{},
		&ProcessorVariable{},
		&ProcessorAssertion{},
		&ProcessorData{},
		&ProcessorCookie{},
		&ProcessorExtractor{},
		&ProcessorCustomCode{},

		&ScenarioReport{},
		&PlanReport{},
		&ExecLogProcessor{},
		&ExecLogExtractor{},
		&ExecLogCheckpoint{},
		&ExecLogScript{},
		&ExecLogDatabaseOpt{},

		&ComponentSchema{},
		&ComponentSchemaSecurity{},

		&Endpoint{},
		&EndpointPathParam{},
		&EndpointInterfaceRequestBody{},
		&EndpointInterfaceRequestBodyItem{},
		&EndpointInterfaceResponseBodyItem{},
		&EndpointInterfaceResponseBodyHeader{},
		&EndpointInterfaceResponseBody{},
		&EndpointInterface{},
		&EndpointCase{},
		&EndpointCaseAlternative{},
		&EndpointCaseAlternativeFactor{},
		&EndpointInterfaceParam{},
		&EndpointInterfaceCookie{},
		&EndpointInterfaceHeader{},
		&EndpointDocument{},
		&EndpointSnapshot{},
		&EndpointTag{},
		&EndpointTagRel{},
		&EndpointInterfaceGlobalParam{},

		&Serve{},
		&ServeServer{},
		&ServeVersion{},
		&EndpointVersion{},
		&ServeEndpointVersion{},
		&SummaryBugs{},
		&SummaryDetails{},
		&SummaryProjectUserRanking{},
		&EnvironmentParam{},
		&Message{},
		&MessageRead{},
		&DebugInvoke{},
		&ProjectPerm{},
		&ProjectRolePerm{},
		&ProjectRoleMenu{},
		&ProjectMenu{},
		&ProjectRecentlyVisited{},
		&ProjectMemberAudit{},
		&DatabaseConn{},
		&LlmTool{},
		&DebugConditionDatabaseOpt{},

		&SwaggerSync{},
		&ProjectMockSetting{},

		&EndpointMockScript{},
		&EndpointMockExpect{},
		&EndpointMockExpectRequest{},
		&EndpointMockExpectResponse{},
		&EndpointMockExpectResponseHeader{},

		&DebugConditionResponseDefine{},
		&ExecLogResponseDefine{},

		&SysConfig{},
		&Jslib{},
		&SysAgent{},
		&ProjectUserServer{},
		&ThirdPartySync{},

		&DebugInterfaceGlobalParam{},
		&DebugInterfaceCookie{},

		&ProjectProductRel{},
		&ProjectSpaceRel{},

		&EndpointFavorite{},

		&ProjectCron{},
		&CronConfigLecang{},
		&ProjectEngineeringRel{},

		&AiMeasurement{},
		&AiMetrics{},
		&AiMetricsAnswerRelevancy{},
		&AiModel{},
		&AiTemplate{},
	}
)
