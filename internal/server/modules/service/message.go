package service

import (
	"encoding/json"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/cron"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/im"
	"github.com/aaronchen2k/deeptest/internal/server/core/cache"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
)

type MessageService struct {
	MessageRepo     *repo.MessageRepo     `inject:""`
	MessageReadRepo *repo.MessageReadRepo `inject:""`
	UserRepo        *repo.UserRepo        `inject:""`
	ProjectRepo     *repo.ProjectRepo     `inject:""`
	ProjectRoleRepo *repo.ProjectRoleRepo `inject:""`
	EndpointRepo    *repo.EndpointRepo    `inject:""`
	BaseRepo        *repo.BaseRepo        `inject:""`
	ProjectService  *ProjectService       `inject:""`
	Cron            *cron.ServerCron      `inject:""`
}

func (s *MessageService) GetScope(userId uint) (scope map[int][]string) {
	return s.MessageRepo.GetScope(userId)
}

func (s *MessageService) Create(req v1.MessageReq) (uint, *_domain.BizErr) {
	return s.MessageRepo.Create(req)
}

func (s *MessageService) Paginate(req v1.MessageReqPaginate, userId uint) (ret _domain.PageData, err error) {
	req.Scope = s.MessageRepo.GetScope(userId)

	ret, err = s.MessageRepo.Paginate(req)

	if err != nil {
		return
	}

	return
}

func (s *MessageService) UnreadCount(userId uint) (count int64, err error) {
	scope := s.MessageRepo.GetScope(userId)
	req := v1.MessageScope{Scope: scope}

	count, err = s.MessageRepo.GetUnreadCount(req)

	if err != nil {
		return
	}

	return
}

func (s *MessageService) OperateRead(req v1.MessageReadReq) (uint, error) {
	return s.MessageReadRepo.Create(req)
}

func (s *MessageService) GetEndpointMcsData(projectId, endpointId uint) (mcsData im.EnterpriseWechatInfoData, err error) {
	userAccount, err := s.ProjectRepo.GetUsernamesByProjectAndRole(projectId, 0, "")
	if err != nil {
		return
	}

	endpoint, err := s.EndpointRepo.Get(endpointId)
	fmt.Println(endpoint)
	if err != nil {
		return
	}

	mcsData = im.EnterpriseWechatInfoData{
		//Content:     "{\\\"articles\\\":[{\\\"title\\\":\\\"接口变更通知\\\",\\\"description\\\":\\\"接口变更通知\\\",\\\"url\\\":\\\"URL\\\",\\\"appid\\\":\\\"wx123123123123123\\\",\\\"pagepath\\\":\\\"http://localhost:8000/Lecang8/IM?shareInfo=%7B%22endpointId%22%3A59047%2C%22selectedCategoryId%22%3A%22%22%7D\\\"}]}", //TODO
		UserAccount: userAccount,
		ImAppid:     config.CONFIG.Mcs.ImAppid,
		MsgType:     1,
	}

	return
}

func (s *MessageService) GetAuditProjectResultMcsData(auditId uint) (mcsData im.EnterpriseWechatInfoData, err error) {
	auditData, err := s.ProjectRepo.GetAudit(auditId)
	if err != nil {
		return
	}

	applyProject, err := s.ProjectRepo.Get(auditData.ProjectId)
	if err != nil {
		return
	}

	applyProjectRole, err := s.ProjectRoleRepo.FindByName(auditData.ProjectRoleName)
	if err != nil {
		return
	}

	applyUser, err := s.UserRepo.GetByUserId(auditData.ApplyUserId)
	if err != nil {
		return
	}

	auditRes := "被拒绝"
	if auditData.Status == consts.Agreed {
		auditRes = "通过"
	}

	host, _ := cache.GetCacheString("host")
	articles := make([]im.EnterpriseWechatInfoContentArticles, 0)
	articles = append(articles, im.EnterpriseWechatInfoContentArticles{
		Title:       "申请加入项目审批结果通知",
		Description: fmt.Sprintf("您申请加入的项目%s【%s角色】审批%s", applyProject.Name, applyProjectRole.Name, auditRes),
		Pagepath:    fmt.Sprintf("%s/sys-setting/user-manage", host),
	})
	content := im.EnterpriseWechatInfoContent{
		Articles: articles,
	}
	contentByte, _ := json.Marshal(content)

	mcsData = im.EnterpriseWechatInfoData{
		Content:     string(contentByte),
		UserAccount: []string{applyUser.Username},
		ImAppid:     config.CONFIG.Mcs.ImAppid,
		MsgType:     1,
	}

	return
}

func (s *MessageService) ReceiveMcsApprovalResult(res v1.McsApprovalResData) (err error) {
	message, err := s.MessageRepo.GetByMcsMessageId(res.InstanceId)
	if err != nil {
		return
	}

	err = s.MessageRepo.UpdateSendStatusByMcsMessageId(res.InstanceId, s.TransferToSendStatus(res.Status))
	if err != nil {
		return
	}

	if res.Status == 5 {
		return
	}

	approveUserName := res.ApproveUser[0]
	approveUser, err := s.UserRepo.GetByUserName(approveUserName)
	if err != nil {
		return
	}

	if message.MessageSource == consts.MessageSourceJoinProject {
		status := consts.Refused
		if res.Status == 2 {
			status = consts.Agreed
		}

		err = s.ProjectService.Audit(message.BusinessId, approveUser.ID, status)
	}

	return
}

func (s *MessageService) TransferToSendStatus(status int) (sendStatus consts.MessageSendStatus) {
	switch status {
	case 1:
		sendStatus = consts.MessageApprovalReject
	case 2:
		sendStatus = consts.MessageApprovalAgreed
	case 5:
		sendStatus = consts.MessageApprovalInProgress
	default:
		sendStatus = consts.MessageApprovalReject
	}

	return
}
