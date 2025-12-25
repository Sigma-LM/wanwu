package request

// note: app_key命名的相关文件是针对应用生成的key（老版本），而api_key是针对用户生成的openapikey。
type GenAppKeyRequest struct {
	AppId   string `json:"appId" validate:"required"`   // 应用id
	AppType string `json:"appType" validate:"required"` // 应用类型
}

func (g GenAppKeyRequest) Check() error {
	return nil
}

type DelAppKeyRequest struct {
	ApiId string `json:"apiId" validate:"required"` // ApiID
}

func (d DelAppKeyRequest) Check() error {
	return nil
}

type GetAppKeyListRequest struct {
	AppId   string `form:"appId" json:"appId" validate:"required"`     // 应用id
	AppType string `form:"appType" json:"appType" validate:"required"` // 应用类型
}

func (g GetAppKeyListRequest) Check() error {
	return nil
}
