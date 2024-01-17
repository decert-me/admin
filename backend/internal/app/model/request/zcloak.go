package request

type GenerateCardInfoRequest struct {
	TokenId string `json:"token_id"`
	Answer  string `json:"answer" binding:"required"`
	Uri     string `json:"uri"`
}
