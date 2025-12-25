package model

// note: 用户api key管理表

type OpenApiKey struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt int64  `gorm:"autoCreateTime:milli;index:idx_api_key_created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
	// 组织ID
	OrgID string `gorm:"index:idx_api_key_org_id"`
	// 用户ID
	UserID string `gorm:"index:idx_api_key_user_id"`
	// Key
	Key string `gorm:"index:idx_api_key_key"`
	// 描述
	Description string `gorm:"index:idx_api_key_description"`
	// 名称
	Name string `gorm:"index:idx_api_key_name"`
	// 是否启用
	Status bool `gorm:"index:idx_api_key_status"`
	// 到期时间
	ExpiredAt int64 `gorm:"index:idx_api_key_expired_at"`
}
