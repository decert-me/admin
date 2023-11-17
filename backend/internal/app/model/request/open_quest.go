package request

type GetOpenQuestPermListRequest struct {
	PageInfo
}

type AddOpenQuestPermRequest struct {
	Address string `json:"address"`
}

type DeleteOpenQuestPermRequest struct {
	Address string `json:"address"`
}
