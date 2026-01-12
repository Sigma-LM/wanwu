package model

type ModelExperienceDialog struct {
	ID           uint32 `gorm:"primary_key;auto_increment;not null;"`
	ModelId      string `gorm:"column:model_id;index:idx_idx_model_experience_model;type:varchar(100);comment:模型 ID"`
	SessionId    string `gorm:"column:session_id;uniqueIndex:idx_model_experience;type:varchar(100);comment:会话 ID"`
	Title        string `gorm:"column:title;type:varchar(512);comment:对话标题"`
	ModelSetting string `gorm:"column:model_setting;type:longtext;comment:模型参数配置"`
	PublicModel
}
