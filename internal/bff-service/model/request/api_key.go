package request

type CreateAPIKeyRequest struct {
	Name      string `json:"name" validate:"required"`
	Desc      string `json:"desc"`
	ExpiredAt string `json:"expiredAt"` // 格式 yyyy-mm-dd
}

func (c *CreateAPIKeyRequest) Check() error {
	return nil
}

type DeleteAPIKeyRequest struct {
	KeyID string `json:"keyId"`
}

func (d *DeleteAPIKeyRequest) Check() error {
	return nil
}

type UpdateAPIKeyRequest struct {
	KeyID     string `json:"keyId"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	ExpiredAt string `json:"expiredAt"`
}

func (u *UpdateAPIKeyRequest) Check() error {
	return nil
}

type UpdateAPIKeyStatusRequest struct {
	KeyID  string `json:"keyId"`
	Status bool   `json:"status"` // 状态 false-禁用 true-启用
}

func (u *UpdateAPIKeyStatusRequest) Check() error {
	return nil
}
