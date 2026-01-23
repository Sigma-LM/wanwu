package orm

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
)

// CreateKnowledgeExternalAPI 创建外部知识库API
func CreateKnowledgeExternalAPI(ctx context.Context, externalAPI *model.KnowledgeExternalAPI) error {
	return db.GetHandle(ctx).Create(externalAPI).Error
}

// UpdateKnowledgeExternalAPI 更新外部知识库API
func UpdateKnowledgeExternalAPI(ctx context.Context, externalAPI *model.KnowledgeExternalAPI) error {
	return db.GetHandle(ctx).Model(&model.KnowledgeExternalAPI{}).
		Where("external_api_id = ?", externalAPI.ExternalAPIId).
		Updates(externalAPI).Error
}

// DeleteKnowledgeExternalAPI 删除外部知识库API
func DeleteKnowledgeExternalAPI(ctx context.Context, externalAPIId string) error {
	return db.GetHandle(ctx).Model(&model.KnowledgeExternalAPI{}).
		Where("external_api_id = ?", externalAPIId).
		Delete(&model.KnowledgeExternalAPI{}).Error
}

// GetKnowledgeExternalAPIList 获取外部知识库API列表
func GetKnowledgeExternalAPIList(ctx context.Context, userId, orgId string, externalAPIIds []string) ([]*model.KnowledgeExternalAPI, error) {
	var externalAPIList []*model.KnowledgeExternalAPI
	err := sqlopt.SQLOptions(sqlopt.WithPermit(orgId, userId),
		sqlopt.WithExternalAPIIdList(externalAPIIds)).Apply(db.GetHandle(ctx), &model.KnowledgeExternalAPI{}).
		Order("create_at desc").Find(&externalAPIList).Error
	if err != nil {
		log.Errorf("GetKnowledgeExternalAPIList err: %v", err)
		return nil, err
	}
	return externalAPIList, nil
}

// GetKnowledgeExternalAPIInfo 获取外部知识库API详情
func GetKnowledgeExternalAPIInfo(ctx context.Context, userId, orgId string, externalAPIId string) (*model.KnowledgeExternalAPI, error) {
	var externalAPIInfo *model.KnowledgeExternalAPI
	err := sqlopt.SQLOptions(sqlopt.WithExternalAPIId(externalAPIId),
		sqlopt.WithPermit(orgId, userId)).Apply(db.GetHandle(ctx), &model.KnowledgeExternalAPI{}).First(&externalAPIInfo).Error
	if err != nil {
		log.Errorf("GetKnowledgeExternalAPIInfo err: %v", err)
		return nil, err
	}
	return externalAPIInfo, nil
}
