package run

import (
	_ "embed"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// SessionRunner is used to run testcase and its steps.
// each testcase has its own SessionRunner instance and share session variables.
type SessionRunner struct {
	*ScenarioRunner
	sessionVariables map[string]interface{}

	startTime        time.Time                  // record start time of the testcase
	summary          *TestScenarioSummary       // record test case summary
	wsConnMap        map[string]*websocket.Conn // save all websocket connections
	pongResponseChan chan string                // channel used to receive pong response message
}

func (r *SessionRunner) resetSession() {
	log.Info().Msg("reset session runner")
	r.sessionVariables = make(map[string]interface{})
	r.startTime = time.Now()
	r.summary = newSummary()
	r.wsConnMap = make(map[string]*websocket.Conn)
	r.pongResponseChan = make(chan string, 1)
}

func (r *SessionRunner) GetParser() *Parser {
	return r.parser
}

func (r *SessionRunner) GetConfig() *TConfig {
	return r.parsedConfig
}

func (r *SessionRunner) HTTPStatOn() bool {
	return r.hrpRunner.httpStatOn
}

func (r *SessionRunner) LogOn() bool {
	return r.hrpRunner.requestsLogOn
}

// Start runs the test steps in sequential order.
// givenVars is used for data driven
func (r *SessionRunner) Start(givenVars map[string]interface{}) error {
	log.Info().Str("testcase", "").Msg("run testcase start")

	// reset session runner
	r.resetSession()

	// update config variables with given variables
	r.updateSessionVariables(givenVars)

	// run step in sequential order
	for _, stage := range r.testScenario.TestStages {
		// parse step name
		parsedName, err := r.parser.ParseString(stage.Name(), r.sessionVariables)
		if err != nil {
			parsedName = stage.Name()
		}
		stepName := convertString(parsedName)
		log.Info().Str("step", stepName).
			Str("type", string(stage.Category())).Msg("run step start")

		// run stage like processor or interface
		stepResult, err := stage.Run(r)

		stepResult.Name = stepName
		if err != nil {
			log.Error().
				Str("step", stepResult.Name).
				Str("type", string(stepResult.StageType)).
				Bool("success", false).
				Msg("run step end")

			if r.hrpRunner.failfast {
				return errors.Wrap(err, "abort running due to failfast setting")
			}
		}

		// update extracted variables
		for k, v := range stepResult.ExportVars {
			r.sessionVariables[k] = v
		}

		log.Info().
			Str("step", stepResult.Name).
			Str("type", string(stepResult.StageType)).
			Bool("success", stepResult.Success).
			Interface("exportVars", stepResult.ExportVars).
			Msg("run step end")
	}

	defer func() {
	}()

	log.Info().Str("testcase", "").Msg("run testcase end")
	return nil
}

// MergeStepVariables merges step variables with config variables and session variables
func (r *SessionRunner) MergeStepVariables(vars map[string]interface{}) (map[string]interface{}, error) {
	// override variables
	// step variables > session variables (extracted variables from previous steps)
	overrideVars := mergeVariables(vars, r.sessionVariables)
	// step variables > testcase config variables
	overrideVars = mergeVariables(overrideVars, r.parsedConfig.Variables)

	// parse step variables
	parsedVariables, err := r.parser.ParseVariables(overrideVars)
	if err != nil {
		log.Error().Interface("variables", r.parsedConfig.Variables).
			Err(err).Msg("parse step variables failed")
		return nil, err
	}
	return parsedVariables, nil
}

// updateSessionVariables updates session variables with given variables.
// this is used for data driven
func (r *SessionRunner) updateSessionVariables(parameters map[string]interface{}) {
	if len(parameters) == 0 {
		return
	}

	log.Info().Interface("parameters", parameters).Msg("update session variables")
	for k, v := range parameters {
		r.sessionVariables[k] = v
	}
}

func (r *SessionRunner) GetSummary() *TestScenarioSummary {
	caseSummary := r.summary
	caseSummary.Name = r.parsedConfig.Name
	caseSummary.Time.StartAt = r.startTime
	caseSummary.Time.Duration = time.Since(r.startTime).Seconds()

	exportVars := make(map[string]interface{})
	for _, value := range r.parsedConfig.Export {
		exportVars[value] = r.sessionVariables[value]
	}

	caseSummary.InOut.ExportVars = exportVars
	caseSummary.InOut.ConfigVars = r.parsedConfig.Variables

	return caseSummary
}