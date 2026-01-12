package model

type ModelExperienceFile struct {
	ID          uint32 `gorm:"primary_key;auto_increment;not null;"`
	FileName    string `gorm:"column:file_name;type:varchar(100);comment:文件名"`
	FilePath    string `gorm:"column:file_path;type:varchar(100);comment:文件路径"`
	FileExt     string `gorm:"column:file_ext;type:varchar(100);comment:文件扩展名"`
	ExtractText string `gorm:"column:extract_text;type:longtext;comment:提取文本"`
	FileSize    int64  `gorm:"column:file_size;type:bigint;comment:文件大小"`
	PublicModel
}
