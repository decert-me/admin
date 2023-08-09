package response

import (
	"backend/internal/app/model"
)

type SysAuthorityResponse struct {
	Authority model.Authority `json:"authority"`
}

type SysAuthorityCopyResponse struct {
	Authority      model.Authority `json:"authority"`
	OldAuthorityId string          `json:"oldAuthorityId"` // 旧角色ID
}

type AuthorityResponse struct {
	AuthorityId       string `json:"authorityId"`
	AuthoritySourceId []uint `json:"authoritySourceId"`
}
