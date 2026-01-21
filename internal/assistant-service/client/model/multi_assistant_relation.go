package model

type MultiAgentRelation struct {
	Id           uint32 `gorm:"column:id;primary_key;type:bigint(20) auto_increment;not null;comment:'id';" json:"id"`         // Primary Key
	MultiAgentId uint32 `gorm:"column:multi_agent_id;index:idx_multi_agent_id,priority:1;type:bigint(20)" json:"multiAgentId"` // 多智能体id
	Version      string `gorm:"column:version;type:varchar(64);NOT NULL;default:'';comment:多智能体版本" json:"version"`             // 多智能体版本
	AgentId      uint32 `gorm:"column:agent_id;index:idx_multi_agent_id,priority:2;type:bigint(20)" json:"agentId"`            // 单智能体id
	CreatedAt    int64  `gorm:"column:create_at;autoCreateTime:milli;type:bigint(20);not null;" json:"createAt"`               // Create Time
	UpdatedAt    int64  `gorm:"column:update_at;autoUpdateTime:milli;type:bigint(20);not null;" json:"updateAt"`               // Update Time
	UserId       string `gorm:"column:user_id;index:idx_user_id_name,priority:1;type:varchar(64);not null;default:'';" json:"userId"`
	OrgId        string `gorm:"column:org_id;type:varchar(64);not null;default:'';" json:"orgId"`
}

func (MultiAgentRelation) TableName() string {
	return "multi_agent_relation"
}
