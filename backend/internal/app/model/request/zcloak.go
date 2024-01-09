package request

type GenerateCardInfoRequest struct {
	TokenId int64  `json:"token_id"`
	Answer  string `json:"answer" binding:"required"`
	Uri     string `json:"uri"`
}
