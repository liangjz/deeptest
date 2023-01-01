package service

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	queryHelper "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/query"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/antchfx/jsonquery"
	"strings"
)

type ParserJsonService struct {
	XPathService      *XPathService      `inject:""`
	ParserRegxService *ParserRegxService `inject:""`
}

func (s *ParserJsonService) ParseJson(req *v1.ParserRequest) (ret v1.ParserResponse, err error) {
	docJson := s.updateJsonElem(req.DocContent, req.SelectContent, req.StartLine, req.EndLine,
		req.StartColumn, req.EndColumn)

	req.DocContent = docJson

	elem := s.getJsonSelectedElem(req.DocContent)

	exprType := "xpath"
	expr, _ := s.XPathService.GetJsonXPath(elem, req.SelectContent, true)
	if expr != "" {
		expr = expr + "/" + req.SelectContent

		result := queryHelper.JsonQuery(req.DocContent, expr)
		fmt.Printf("%s: %v", expr, result)
	} else {
		expr, _ = s.ParserRegxService.getRegxExpr(req.DocContent, req.SelectContent,
			req.StartLine, req.StartColumn,
			req.EndLine, req.EndColumn)
		exprType = "regx"
	}

	ret = v1.ParserResponse{
		SelectionType: consts.NodeProp,
		Expr:          expr,
		ExprType:      exprType,
	}

	return
}

func (s *ParserJsonService) updateJsonElem(docJson, selectContent string,
	startLine, endLine, startColumn, endColumn int) (ret string) {
	lines := strings.Split(docJson, "\n")

	line := []rune(lines[startLine])

	newStr := fmt.Sprintf("%s-%s", consts.DeepestKey, selectContent)
	newLine := string(line[:startColumn]) + newStr + string(line[endColumn:])

	lines[startLine] = newLine

	ret = strings.Join(lines, "\n")
	return
}

func (s *ParserJsonService) getJsonSelectedElem(docJson string) (ret *jsonquery.Node) {
	doc, err := jsonquery.Parse(strings.NewReader(docJson))
	if err != nil {
		return
	}

	expr := fmt.Sprintf("//*[contains(.,'%s')]", consts.DeepestKey)
	ret, err = jsonquery.Query(doc, expr)

	return
}

func (s *ParserJsonService) queryElem(docJson, xpath string) (ret *jsonquery.Node) {
	doc, err := jsonquery.Parse(strings.NewReader(docJson))
	if err != nil {
		return
	}

	expr := fmt.Sprintf(xpath)
	ret, err = jsonquery.Query(doc, expr)

	return
}