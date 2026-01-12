package model

import (
	"context"

	"github.com/UnicomAI/wanwu/pkg/util"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	"github.com/UnicomAI/wanwu/internal/model-service/client/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) SaveModelExperienceDialog(ctx context.Context, req *model_service.ModelExperienceDialogReq) (*model_service.ModelExperienceDialog, error) {
	modelExperienceDialog, err := s.cli.SaveModelExperienceDialog(ctx, &model.ModelExperienceDialog{
		SessionId:    req.SessionId,
		ModelId:      req.ModelId,
		Title:        req.Title,
		ModelSetting: req.ModelSetting,
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})

	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceDialog, err)
	}

	res := &model_service.ModelExperienceDialog{
		Id:           util.Int2Str(modelExperienceDialog.ID),
		SessionId:    modelExperienceDialog.SessionId,
		ModelId:      modelExperienceDialog.ModelId,
		Title:        modelExperienceDialog.Title,
		ModelSetting: modelExperienceDialog.ModelSetting,
		OrgId:        req.OrgId,
		UserId:       req.UserId,
	}
	return res, nil
}

func (s *Service) GetModelExperienceDialogs(ctx context.Context, req *model_service.GetModelExperienceDialogReq) (*model_service.ModelExperienceDialogs, error) {
	modelExperienceDialogs, err := s.cli.ListExperienceDialogs(ctx, &model.ModelExperienceDialog{
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})

	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceDialog, err)
	}

	var res []*model_service.ModelExperienceDialog
	for _, modelExperienceDialog := range modelExperienceDialogs {
		res = append(res, &model_service.ModelExperienceDialog{
			Id:           util.Int2Str(modelExperienceDialog.ID),
			SessionId:    modelExperienceDialog.SessionId,
			ModelId:      modelExperienceDialog.ModelId,
			Title:        modelExperienceDialog.Title,
			ModelSetting: modelExperienceDialog.ModelSetting,
			CreatedAt:    modelExperienceDialog.CreatedAt,
		})

	}
	return &model_service.ModelExperienceDialogs{
		ModelExperienceDialog: res,
		Total:                 int64(len(res)),
	}, nil
}

func (s *Service) GetModelExperienceDialog(ctx context.Context, req *model_service.ModelExperienceDialogIdReq) (*model_service.ModelExperienceDialog, error) {
	modelExperienceDialog, err := s.cli.GetModelExperienceDialog(ctx, &model.ModelExperienceDialog{
		ID: util.MustU32(req.ModelExperienceId),
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})

	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceDialog, err)
	}

	res := &model_service.ModelExperienceDialog{
		Id:           util.Int2Str(modelExperienceDialog.ID),
		SessionId:    modelExperienceDialog.SessionId,
		ModelId:      modelExperienceDialog.ModelId,
		Title:        modelExperienceDialog.Title,
		ModelSetting: modelExperienceDialog.ModelSetting,
		OrgId:        req.OrgId,
		UserId:       req.UserId,
	}
	return res, nil
}

func (s *Service) DeleteModelExperienceDialog(ctx context.Context, req *model_service.ModelExperienceDialogIdReq) (*emptypb.Empty, error) {
	err := s.cli.DeleteModelExperienceDialog(ctx, &model.ModelExperienceDialog{
		ID: util.MustU32(req.ModelExperienceId),
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})
	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceDialog, err)
	}
	return nil, nil
}

func (s *Service) SaveModelExperienceDialogRecord(ctx context.Context, req *model_service.ModelExperienceDialogRecordReq) (*model_service.ModelExperienceDialogRecord, error) {
	modelExperienceDialogRecord, err := s.cli.SaveModelExperienceDialogRecord(ctx, &model.ModelExperienceDialogRecord{
		ModelExperienceID: util.MustU32(req.ModelExperienceId),
		ModelId:           req.ModelId,
		OriginalContent:   req.OriginalContent,
		HandledContent:    req.HandledContent,
		ReasoningContent:  req.ReasoningContent,
		Role:              req.Role,
		ParentID:          util.MustU32(req.ParentID),
		FileIdList:        req.FileIdList,

		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})
	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceDialogRecord, err)
	}
	res := &model_service.ModelExperienceDialogRecord{
		Id:                util.Int2Str(modelExperienceDialogRecord.ID),
		ModelExperienceId: util.Int2Str(modelExperienceDialogRecord.ModelExperienceID),
		ModelId:           modelExperienceDialogRecord.ModelId,
		OriginalContent:   modelExperienceDialogRecord.OriginalContent,
		HandledContent:    modelExperienceDialogRecord.HandledContent,
		ReasoningContent:  modelExperienceDialogRecord.ReasoningContent,
		Role:              modelExperienceDialogRecord.Role,
	}
	return res, nil
}

func (s *Service) GetModelExperienceDialogRecords(ctx context.Context, req *model_service.GetModelExperienceDialogRecordReq) (*model_service.ModelExperienceDialogRecords, error) {
	modelExperienceDialogRecords, err := s.cli.GetModelExperienceDialogRecord(ctx, util.MustU32(req.ModelExperienceId), req.SessionId, req.ModelId)
	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceDialogRecord, err)
	}
	var res []*model_service.ModelExperienceDialogRecord
	for _, modelExperienceDialogRecord := range modelExperienceDialogRecords {
		res = append(res, &model_service.ModelExperienceDialogRecord{
			Id:                util.Int2Str(modelExperienceDialogRecord.ID),
			ModelExperienceId: util.Int2Str(modelExperienceDialogRecord.ModelExperienceID),
			ModelId:           modelExperienceDialogRecord.ModelId,
			OriginalContent:   modelExperienceDialogRecord.OriginalContent,
			HandledContent:    modelExperienceDialogRecord.HandledContent,
			ReasoningContent:  modelExperienceDialogRecord.ReasoningContent,
			Role:              modelExperienceDialogRecord.Role,
			ParentID:          util.Int2Str(modelExperienceDialogRecord.ParentID),
			FileIdList:        modelExperienceDialogRecord.FileIdList,
		})
	}
	return &model_service.ModelExperienceDialogRecords{
		Record: res,
		Total:  int64(len(res)),
	}, nil
}

func (s *Service) SaveModelExperienceFile(ctx context.Context, req *model_service.ModelExperienceFileReq) (*model_service.ModelExperienceFile, error) {
	modelExperienceDialog, err := s.cli.SaveModelExperienceFile(ctx, &model.ModelExperienceFile{
		FileName:    req.FileName,
		FilePath:    req.FilePath,
		FileExt:     req.FileExt,
		ExtractText: req.ExtractText,
		FileSize:    req.FileSize,
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})

	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceFile, err)
	}

	res := &model_service.ModelExperienceFile{
		Id:          util.Int2Str(modelExperienceDialog.ID),
		FileName:    modelExperienceDialog.FileName,
		FilePath:    modelExperienceDialog.FilePath,
		FileExt:     modelExperienceDialog.FileExt,
		ExtractText: modelExperienceDialog.ExtractText,
		FileSize:    modelExperienceDialog.FileSize,
		OrgId:       req.OrgId,
		UserId:      req.UserId,
	}
	return res, nil
}

func (s *Service) GetModelExperienceFilesByIds(ctx context.Context, getModelExperienceFilesReq *model_service.GetModelExperienceFilesByIdsReq) (*model_service.ModelExperienceFiles, error) {
	fileIds := make([]uint32, 0, len(getModelExperienceFilesReq.FileIds))
	for _, fileId := range getModelExperienceFilesReq.FileIds {
		fileIds = append(fileIds, util.MustU32(fileId))
	}
	modelExperienceFiles, err := s.cli.GetModelExperienceFilesByIds(ctx, fileIds)
	if err != nil {
		return nil, errStatus(errs.Code_ModelExperienceFileByIds, err)
	}

	// Convert internal model to proto-generated model
	var res []*model_service.ModelExperienceFile
	for _, modelExperienceFile := range modelExperienceFiles {
		res = append(res, &model_service.ModelExperienceFile{
			Id:          util.Int2Str(modelExperienceFile.ID),
			FileName:    modelExperienceFile.FileName,
			FilePath:    modelExperienceFile.FilePath,
			FileExt:     modelExperienceFile.FileExt,
			ExtractText: modelExperienceFile.ExtractText,
			FileSize:    modelExperienceFile.FileSize,
			OrgId:       modelExperienceFile.OrgID,
			UserId:      modelExperienceFile.UserID,
		})
	}
	return &model_service.ModelExperienceFiles{
		Files: res,
		Total: int64(len(res)),
	}, nil
}
