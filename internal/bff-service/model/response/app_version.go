package response

type AppVersionInfo struct {
	Version     string `json:"version"`
	Desc        string `json:"desc"`
	CreatedAt   string `json:"createdAt"`
	PublishType string `json:"publishType"`
}
