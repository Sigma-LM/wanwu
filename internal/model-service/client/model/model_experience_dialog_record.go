package model

type ModelExperienceDialogRecord struct {
	ID                uint32 `gorm:"primary_key;auto_increment;not null;"`
	ModelExperienceID uint32 `gorm:"column:model_experience_id;index:idx_model_experience_dialog_model_record_experience_id;type:int;comment:模型体验ID"`
	SessionId         string `gorm:"column:SessionId;idx:idx_model_experience;type:varchar(100);comment:会话ID"`
	ModelId           string `gorm:"column:model_id;index:idx_idx_model_experience_model;type:varchar(100);comment:模型 ID"`
	OriginalContent   string `gorm:"column:original_prompt;type:longtext;comment:原始内容"`
	HandledContent    string `gorm:"column:handled_prompt;type:longtext;comment:处理后内容"`
	ReasoningContent  string `gorm:"column:reasoning_prompt;type:longtext;comment:思考过程"`
	Role              string `gorm:"column:role;type:varchar(100);comment:角色"`
	ParentID          uint32 `gorm:"column:parent_id;index:idx_model_experience_dialog_record_parent_id;type:int;"`
	FileIdList        string `gorm:"column:file_id_list;type:varchar(512);comment:文件 ID 列表"`
	PublicModel
}
