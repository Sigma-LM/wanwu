package orm

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	model_client "github.com/UnicomAI/wanwu/internal/model-service/client/model"
	"github.com/UnicomAI/wanwu/internal/model-service/client/orm/sqlopt"
	"gorm.io/gorm/clause"
)

func (c *Client) ListExperienceDialogs(ctx context.Context, req *model_client.ModelExperienceDialog) ([]*model_client.ModelExperienceDialog, *errs.Status) {
	var dialogs []*model_client.ModelExperienceDialog
	if err := sqlopt.SQLOptions(
		sqlopt.WithOrgID(req.OrgID),
		sqlopt.WithUserID(req.UserID),
	).Apply(c.db.WithContext(ctx)).Order("created_at desc").Find(&dialogs).Error; err != nil {
		return nil, toErrStatus("get_experience_dialogs_err", err.Error())
	}
	return dialogs, nil
}

func (c *Client) SaveModelExperienceDialog(ctx context.Context, experience *model_client.ModelExperienceDialog) (*model_client.ModelExperienceDialog, *errs.Status) {
	// 创建更新字段的映射
	updates := map[string]interface{}{
		"model_id":      experience.ModelId,
		"model_setting": experience.ModelSetting,
		"org_id":        experience.OrgID,
		"user_id":       experience.UserID,
	}

	// 只有当title不为空时才更新title字段
	if experience.Title != "" {
		updates["title"] = experience.Title
	}

	if err := c.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_id"}},
		DoUpdates: clause.Assignments(updates),
	}).Create(experience).Error; err != nil {
		return nil, toErrStatus("create_experience_err", err.Error())
	}
	return experience, nil
}

func (c *Client) GetModelExperienceDialog(ctx context.Context, tab *model_client.ModelExperienceDialog) (*model_client.ModelExperienceDialog, *errs.Status) {
	info := &model_client.ModelExperienceDialog{}
	if err := sqlopt.SQLOptions(
		sqlopt.WithID(tab.ID),
		sqlopt.WithOrgID(tab.OrgID),
		sqlopt.WithUserID(tab.UserID),
	).Apply(c.db).WithContext(ctx).First(info).Error; err != nil {
		return nil, toErrStatus("model_experience_dialog_get_err", err.Error())
	}
	return info, nil
}

func (c *Client) DeleteModelExperienceDialog(ctx context.Context, experience *model_client.ModelExperienceDialog) *errs.Status {
	var existing model_client.ModelExperienceDialog
	if err := sqlopt.SQLOptions(
		sqlopt.WithID(experience.ID),
		sqlopt.WithOrgID(experience.OrgID),
		sqlopt.WithUserID(experience.UserID),
	).Apply(c.db).WithContext(ctx).Select("id").First(&existing).Error; err != nil {
		return toErrStatus("delete_model_experience_err", err.Error())
	}
	if err := c.db.WithContext(ctx).Delete(existing).Error; err != nil {
		return toErrStatus("delete_model_experience_err", err.Error())
	}
	return nil
}

func (c *Client) SaveModelExperienceDialogRecord(ctx context.Context, record *model_client.ModelExperienceDialogRecord) (*model_client.ModelExperienceDialogRecord, *errs.Status) {
	if err := c.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, toErrStatus("create_experience_err", err.Error())
	}
	return record, nil
}

func (c *Client) SaveModelExperienceFile(ctx context.Context, file *model_client.ModelExperienceFile) (*model_client.ModelExperienceFile, *errs.Status) {
	if err := c.db.WithContext(ctx).Create(file).Error; err != nil {
		return nil, toErrStatus("create_experience_err", err.Error())
	}
	return file, nil
}

func (c *Client) GetModelExperienceFilesByIds(ctx context.Context, fileIds []uint32) ([]*model_client.ModelExperienceFile, *errs.Status) {
	var files []*model_client.ModelExperienceFile
	if err := c.db.WithContext(ctx).Where("id IN ?", fileIds).Find(&files).Error; err != nil {
		return nil, toErrStatus("get_experience_files_err", err.Error())
	}
	return files, nil
}

func (c *Client) GetModelExperienceDialogRecord(ctx context.Context, modelExperienceId uint32, sessionId, modelId string) ([]*model_client.ModelExperienceDialogRecord, *errs.Status) {
	var records []*model_client.ModelExperienceDialogRecord
	query := c.db.WithContext(ctx)
	if modelExperienceId != 0 {
		query = query.Where("model_experience_id = ?", modelExperienceId)
	} else {
		query = query.Where("session_id = ?", sessionId)
	}
	// 只有当model不为空时，才添加model条件
	if modelId != "" {
		query = query.Where("model_Id = ?", modelId)
	}
	if err := query.Order("id asc").Find(&records).Error; err != nil {
		return nil, toErrStatus("get_experience_record_err", err.Error())
	}
	return records, nil
}
