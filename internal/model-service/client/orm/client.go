package orm

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/model-service/client/model"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(ctx context.Context, db *gorm.DB) (*Client, error) {
	// auto migrate
	if err := db.AutoMigrate(
		model.ModelImported{},
	); err != nil {
		return nil, err
	}

	// 数据初始化处理
	if err := initModelImportedProviderName(db); err != nil {
		return nil, err
	}
	if err := initModelUUID(db); err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}

func toErrStatus(key string, args ...string) *err_code.Status {
	return &err_code.Status{
		TextKey: key,
		Args:    args,
	}
}

const (
	providerHuoshanOriginal = "Huoshan" // 原始值
	providerHuoshanTarget   = "HuoShan" // 目标值
)

// initModelImportedProviderName 更新数据库中ModelImported表的provider字段值
func initModelImportedProviderName(dbClient *gorm.DB) error {
	err := dbClient.Model(&model.ModelImported{}).
		Where("provider = ?", providerHuoshanOriginal).
		Update("provider", providerHuoshanTarget).Error
	if err != nil {
		return err
	}
	return nil
}

// initModelUUID 批量更新数据库中ModelImported表的uuid字段值
func initModelUUID(dbClient *gorm.DB) error {
	const batchSize = 100

	for {
		var ids []uint32
		if err := dbClient.Model(&model.ModelImported{}).Select("id").Where("uuid = ? OR uuid IS NULL", "").Limit(batchSize).Find(&ids).Error; err != nil {
			return err
		}

		if len(ids) == 0 {
			break
		}

		caseWhen := "CASE id "
		var args []interface{}
		for _, id := range ids {
			caseWhen += "WHEN ? THEN ? "
			args = append(args, id, util.NewID())
		}
		caseWhen += "END"

		if err := dbClient.Model(&model.ModelImported{}).
			Where("id IN ?", ids).
			UpdateColumn("uuid", gorm.Expr(caseWhen, args...)).Error; err != nil {
			log.Errorf("init model uuid batch update error: %v", err)
			return err
		}
	}

	return nil
}
