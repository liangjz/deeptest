package main

import (
	"flag"
	"github.com/aaronchen2k/deeptest/cmd/agent/agent"
	v1 "github.com/aaronchen2k/deeptest/cmd/agent/v1"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/cron"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/websocket"
	_consts "github.com/aaronchen2k/deeptest/pkg/consts"
	"github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	AppVersion string
	BuildTime  string
	GoVersion  string
	GitHash    string

	flagSet *flag.FlagSet
)

// @title DeepTest代理API文档
// @version 1.0
// @contact.name API Support
// @contact.url https://github.com/aaronchen2k/deeptest/issues
// @contact.email 462626@qq.com
func main() {
	flagSet = flag.NewFlagSet("deeptest", flag.ContinueOnError)

	flagSet.IntVar(&consts.Port, "p", 0, "")
	flagSet.BoolVar(&_consts.Verbose, "verbose", false, "")

	websocketHelper.InitMq()

	server := agent.Init()
	if server == nil {
		return
	}

	injectModule(server)
	server.Start()
}

func injectModule(ws *agent.WebServer) {
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	cron := cron.NewServerCron()
	cron.Init()
	indexModule := v1.NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: cron},
		&inject.Object{Value: indexModule},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}

	ws.AddModule(indexModule.Party())

	_logUtils.Infof("start agent")
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}