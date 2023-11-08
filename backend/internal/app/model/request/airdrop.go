package request

type RunAirdropReq struct {
	App string `json:"app"`
}

type GetAirdropListReq struct {
	PageInfo
	App string `json:"app"`
}
