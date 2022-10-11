package agentExec

import (
	"crypto/tls"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/kataras/iris/v12/websocket"
	"golang.org/x/net/http2"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Session struct {
	ScenarioId uint
	Name       string

	HttpClient  *http.Client
	Http2Client *http.Client
	Failfast    bool

	RootProcessor *Processor
	Report        *Report

	WsMsg *websocket.Message
}

func NewSession(root *Processor, variables []domain.Variable, failfast bool, wsMsg *websocket.Message) (ret *Session) {
	ImportVariables(root.ID, variables, consts.Global)

	session := Session{
		ScenarioId:    root.ScenarioId,
		Name:          root.Name,
		RootProcessor: root,
		Failfast:      failfast,
		WsMsg:         wsMsg,
	}

	jar, _ := cookiejar.New(nil)
	session.HttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar:     jar, // insert response cookies into request
		Timeout: 120 * time.Second,
	}

	session.Http2Client = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 120 * time.Second,
	}

	ret = &session

	return
}

func (s *Session) Run() {
	exec.SendStartMsg(s.WsMsg)

	s.RootProcessor.Run(s)

	exec.SendEndMsg(s.WsMsg)
}
