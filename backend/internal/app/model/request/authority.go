package request

type SetDataAuthorityRequest struct {
	AuthorityId       string `json:"authorityId" form:"authorityId"`
	AuthoritySourceId []uint `json:"authoritySourceId"`
}
