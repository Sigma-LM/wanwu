package response

import "github.com/UnicomAI/wanwu/internal/bff-service/model/request"

type Response struct {
	Code int64       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
}

type ListResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserPermission struct {
	OrgPermission    UserOrgPermission `json:"orgPermission"`    // 用户所在组织权限
	Language         Language          `json:"language"`         // 语言
	IsUpdatePassword bool              `json:"isUpdatePassword"` // 是否已更新密码
	Avatar           request.Avatar    `json:"avatar"`           // 用户头像信息
}

type UserOrgPermission struct {
	IsAdmin     bool         `json:"isAdmin"`     // 是否系统内置管理员
	IsSystem    bool         `json:"isSystem"`    // 是否系统视角（此时org.id为空，org.name为"系统"）
	Org         IDName       `json:"org"`         // 组织
	Roles       []IDName     `json:"roles"`       // 角色列表
	Permissions []Permission `json:"permissions"` // 权限列表
}

type Permission struct {
	Perm string `json:"perm"` // 权限
	Name string `json:"name"` // 权限名（对应菜单名）
}

type Select struct {
	Select []IDName `json:"select"`
}

type DocMenu struct {
	Name     string     `json:"name"`     // 目录名称
	Index    string     `json:"index"`    // 目录索引
	Path     string     `json:"path"`     // 目录路径（转码后）
	PathRaw  string     `json:"pathRaw"`  // 目录路径
	Children []*DocMenu `json:"children"` // 目录

	content string
}

func (dm *DocMenu) SetContent(content string) {
	dm.content = content
}

type DocSearchResp struct {
	Title       string             `json:"title"` // 文档名
	ContentList []DocSearchContent `json:"list"`  // 内容列表
}

type DocSearchContent struct {
	Title   string `json:"title"`   // 文档中的子标题
	Content string `json:"content"` // 内容
	Url     string `json:"url"`     // 文档链接
}

type StatisticChart struct {
	TableName string               `json:"tableName"` // 统计表名字
	Lines     []StatisticChartLine `json:"lines"`     // 统计表中线段集合
}

type StatisticChartLine struct {
	LineName string                   `json:"lineName"` // 线段名字
	Items    []StatisticChartLineItem `json:"items"`    // 线段横纵坐标值
}

type StatisticChartLineItem struct {
	Key   string  `json:"key"`
	Value float32 `json:"value"`
}
