package response

// note: app_key命名的相关文件是针对应用生成的key（老版本），而api_key是针对用户生成的全局唯一key。
type AppResponse struct {
	ApiID     string `json:"apiId" `    // ApiID
	ApiKey    string `json:"apiKey"`    // 生成的ApiKey
	CreatedAt string `json:"createdAt"` // 创建ApiKey的时间
}
