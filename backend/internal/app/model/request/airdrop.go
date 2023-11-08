package request

type RunAirdropReq struct {
	App string `json:"app"`
}

type GetAirdropListReq struct {
	PageInfo
	App    string `json:"app"`
	Status uint8  `json:"status" form:"status"` // 状态 0 处理中 1 交易成功 2 交易失败 3 超过解析次数 4 事件匹配失败 5 出现错误
}
