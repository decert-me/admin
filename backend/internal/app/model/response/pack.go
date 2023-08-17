package response

import (
	"time"
)

type PackListResponse struct {
	ID         uint `json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Img        string `json:"img,omitempty"`                                             // 教程封面图
	Label      string `json:"label,omitempty"`                                           // 教程名称
	PackStatus uint8  `gorm:"column:pack_status;default:1" json:"pack_status,omitempty"` // 状态 1 未打包 2 打包成功 3 打包失败
	Branch     string `json:"branch,omitempty"`
	DocPath    string `json:"docPath,omitempty"`
	CommitHash string `json:"commitHash,omitempty"`
}
