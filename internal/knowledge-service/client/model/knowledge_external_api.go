package model

const (
	KnowledgeExternalAPIProviderDify = "dify"
)

type KnowledgeExternalAPI struct {
	Id            uint32 `json:"id" gorm:"primary_key;type:bigint(20) auto_increment;not null;comment:'id';"` // Primary Key
	ExternalAPIId string `gorm:"uniqueIndex:idx_unique_external_api_id;column:external_api_id;type:varchar(64);not null;default:''" json:"externalAPIId"`
	Name          string `gorm:"column:name;type:varchar(256);not null;default:''" json:"name"`
	Description   string `gorm:"column:description;type:text;" json:"description"`
	BaseUrl       string `gorm:"column:base_url;type:varchar(256);not null;default:''" json:"baseUrl"`
	APIKey        string `gorm:"column:api_key;type:varchar(256);not null;default:''" json:"apiKey"`
	Provider      string `gorm:"column:provider;type:varchar(256);not null;default:''" json:"provider"`
	CreatedAt     int64  `gorm:"column:create_at;type:bigint(20);not null;autoCreateTime:milli" json:"createAt"` // Create Time
	UpdatedAt     int64  `gorm:"column:update_at;type:bigint(20);not null;autoUpdateTime:milli" json:"updateAt"` // Update Time
	UserId        string `gorm:"column:user_id;type:varchar(64);not null;default:'';" json:"userId"`
	OrgId         string `gorm:"column:org_id;type:varchar(64);not null;default:''" json:"orgId"`
}

func (KnowledgeExternalAPI) TableName() string {
	return "knowledge_external_api"
}
