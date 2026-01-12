package client

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/model-service/client/model"
)

type IClient interface {
	ImportModel(ctx context.Context, req *model.ModelImported) *errs.Status
	UpdateModel(ctx context.Context, req *model.ModelImported) *errs.Status
	DeleteModel(ctx context.Context, req *model.ModelImported) *errs.Status
	ChangeModelStatus(ctx context.Context, req *model.ModelImported) *errs.Status
	GetModel(ctx context.Context, req *model.ModelImported) (*model.ModelImported, *errs.Status)
	GetModelByUUID(ctx context.Context, uuid string) (*model.ModelImported, *errs.Status)
	GetModelByIds(ctx context.Context, modelIds []uint32) ([]*model.ModelImported, *errs.Status)
	ListModels(ctx context.Context, req *model.ModelImported) ([]*model.ModelImported, *errs.Status)
	ListTypeModels(ctx context.Context, req *model.ModelImported) ([]*model.ModelImported, *errs.Status)

	// Model experience methods
	ListExperienceDialogs(ctx context.Context, req *model.ModelExperienceDialog) ([]*model.ModelExperienceDialog, *errs.Status)
	SaveModelExperienceDialog(ctx context.Context, experience *model.ModelExperienceDialog) (*model.ModelExperienceDialog, *errs.Status)
	GetModelExperienceDialog(ctx context.Context, experience *model.ModelExperienceDialog) (*model.ModelExperienceDialog, *errs.Status)
	DeleteModelExperienceDialog(ctx context.Context, experience *model.ModelExperienceDialog) *errs.Status
	SaveModelExperienceDialogRecord(ctx context.Context, experience *model.ModelExperienceDialogRecord) (*model.ModelExperienceDialogRecord, *errs.Status)
	GetModelExperienceDialogRecord(ctx context.Context, model_experience_id uint32, sessionId, model string) ([]*model.ModelExperienceDialogRecord, *errs.Status)
	SaveModelExperienceFile(ctx context.Context, experience *model.ModelExperienceFile) (*model.ModelExperienceFile, *errs.Status)
	GetModelExperienceFilesByIds(ctx context.Context, fileIds []uint32) ([]*model.ModelExperienceFile, *errs.Status)
}
