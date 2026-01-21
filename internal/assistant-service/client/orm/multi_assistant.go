package orm

import (
	"context"
	"encoding/json"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type MultiAgentDetail struct {
	MultiAgent         *model.Assistant
	MultiAgentSnapshot *model.AssistantSnapshot
	SubAgents          []*model.AssistantSnapshot
}

func (c *Client) GetMultiAssistant(ctx context.Context, multiAssistantID uint32, userID, orgID string, draft bool, version string) (multiAgent *model.Assistant, multiAgentSnapshot *model.AssistantSnapshot, subAgents []*model.AssistantSnapshot, err error) {
	tx := c.db.WithContext(ctx)
	multiAgentDetail, err := searchMultiAssistant(tx, multiAssistantID, userID, orgID, draft, version)
	if err != nil {
		log.Errorf("获取多智能体详情失败 错误(%v) 参数(%v)", err, multiAssistantID)
		return nil, nil, nil, ErrCode(err_code.Code_AssistantErr)
	}
	var multiVersion = ""
	if !draft {
		multiVersion = multiAgentDetail.MultiAgentSnapshot.Version
	}
	multiAgentDetail.SubAgents, err = getSubAssistantList(tx, multiAssistantID, userID, orgID, true, multiVersion)
	if err != nil {
		log.Errorf("获取多智能体详情失败 错误(%v) 参数(%v)", err, multiAssistantID)
		return nil, nil, nil, ErrCode(err_code.Code_AssistantErr)
	}
	return multiAgentDetail.MultiAgent, multiAgentDetail.MultiAgentSnapshot, multiAgentDetail.SubAgents, nil
}

func searchMultiAssistant(tx *gorm.DB, multiAssistantID uint32, userID, orgID string, draft bool, version string) (*MultiAgentDetail, error) {
	multiAgentDetail := &MultiAgentDetail{}
	var err error
	if draft {
		multiAgentDetail.MultiAgent, err = getMultiAssistantDetail(tx, multiAssistantID, userID, orgID)
		if err != nil {
			return nil, err
		}
	} else {
		multiAgentDetail.MultiAgentSnapshot, err = getMultiAssistantSnapshot(tx, multiAssistantID, version)
		if err == nil {
			assistant := &model.Assistant{}
			err = json.Unmarshal([]byte(multiAgentDetail.MultiAgentSnapshot.AssistantInfo), assistant)
			if err != nil {
				return nil, err
			}
			multiAgentDetail.MultiAgent = assistant
		}
	}
	return multiAgentDetail, nil
}

func getMultiAssistantDetail(tx *gorm.DB, multiAssistantID uint32, userID, orgID string) (*model.Assistant, error) {
	var assistant = &model.Assistant{}
	err := sqlopt.SQLOptions(
		sqlopt.WithID(multiAssistantID),
	).Apply(tx).First(assistant).Error
	if err != nil {
		return nil, err
	}
	return assistant, nil
}

func getMultiAssistantSnapshot(tx *gorm.DB, multiAssistantID uint32, version string) (*model.AssistantSnapshot, error) {
	assistantSnapshot := &model.AssistantSnapshot{}
	err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(multiAssistantID),
		sqlopt.WithVersion(version),
	).Apply(tx).Model(&model.AssistantSnapshot{}).
		Order("created_at DESC").
		First(&assistantSnapshot).Error
	if err != nil {
		return nil, err
	}
	return assistantSnapshot, nil
}

func getSubAssistantList(tx *gorm.DB, multiAssistantID uint32, userID, orgID string, needDetail bool, version string) ([]*model.AssistantSnapshot, error) {
	//如果不需要子智能体详情，则只返回子智能体的id
	if !needDetail {
		var subAssistantList []*model.MultiAgentRelation
		err := tx.Select("multi_agent_id", "agent_id").Where("multi_agent_id = ?", multiAssistantID).Where("version = ?", version).Find(&subAssistantList).Error
		if err != nil {
			return nil, err
		}
		return lo.Map(subAssistantList, func(item *model.MultiAgentRelation, index int) *model.AssistantSnapshot {
			return &model.AssistantSnapshot{
				AssistantID: item.AgentId,
			}
		}), nil
	} else {
		var subAssistantList []*model.AssistantSnapshot
		//此查询会查出子智能体的所有版本（sql逻辑简单）,但是如果子智能体版本过多可能有一定内存性能问题
		//如果要一次性查出最新版本则需要很复杂
		err := tx.Model(&model.MultiAgentRelation{}).
			Select("assistant_snapshots.*").
			Joins("LEFT JOIN assistant_snapshots ON multi_agent_relation.agent_id = assistant_snapshots.assistant_id").
			Where("multi_agent_relation.multi_agent_id = ?", multiAssistantID).
			Where("multi_agent_relation.version = ?", version). //版本为空标识draft版本
			Order("assistant_snapshots.created_at DESC").
			Scan(&subAssistantList).Error

		//因为按时间排序了，这里按assistant + version 进行去重取出最新版本
		if len(subAssistantList) > 0 {
			var dataMap = make(map[uint32]bool)
			var retList []*model.AssistantSnapshot
			for i := 0; i < len(subAssistantList); i++ {
				if _, ok := dataMap[subAssistantList[i].AssistantID]; !ok {
					dataMap[subAssistantList[i].AssistantID] = true
					retList = append(retList, subAssistantList[i])
				}
			}
			subAssistantList = retList
		}
		return subAssistantList, err
	}
}
