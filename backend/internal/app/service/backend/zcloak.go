package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"encoding/json"
	"errors"
	"fmt"
	reqV3 "github.com/imroc/req/v3"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"strings"
)

type DidCardRequest struct {
	Receiver string               `json:"receiver"`
	Params   DidCardParamsRequest `json:"params"`
}

type DidCardParamsRequest struct {
	Title       string `json:"Title"`
	ChallengeID string `json:"ChallengeID"`
	Pass        bool   `json:"Pass"`
	Score       int64  `json:"Score"`
	Content     string `json:"Content"`
}

// GenerateCardInfo 生成 card 信息
func GenerateCardInfo(address string, score int64, req request.GenerateCardInfoRequest) (err error) {
	// 获取did 账号
	var did string
	err = global.DB.
		Model(&model.ZcloakDid{}).
		Select("did_address").
		Where("address", address).
		First(&did).Error
	if err != nil || did == "" {
		return errors.New("DIDNotFound")
	}
	// 校验分数正确性
	var quest model.Quest
	err = global.DB.
		Model(&model.Quest{}).
		Where("token_id", req.TokenId).
		First(&quest).Error
	if err != nil {
		return errors.New("TokenIDInvalid")
	}
	// 校验题目
	if req.Uri != "" && req.Uri != quest.Uri {
		return errors.New("QuestUpdate")
	}
	pass := true
	if score == 0 {
		result, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(req.Answer), quest)
		if err != nil {
			global.LOG.Error("AnswerCheck error", zap.Error(err))
			return errors.New("UnexpectedError")
		}
		score = result.UserScore
		score = score / 100
	}

	// 未通过跳过
	if !pass {
		return nil
	}
	// 查询历史 Did 最高分
	var highestScore int64
	if err := global.DB.Model(&model.ZcloakCard{}).
		Select("score").
		Where("did = ? AND quest_id = ?", did, quest.ID).
		Order("score desc").
		First(&highestScore).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			global.LOG.Error("GetHighestScore error", zap.Error(err))
			return errors.New("UnexpectedError")
		}
	}
	// 未达到历史最高分，不保存
	if highestScore >= score {
		return nil
	}
	// 将用户答案写入metadata
	metadata, err := sjson.Set(string(quest.MetaData), "attributes.user_answer", gjson.Parse(req.Answer).Value())
	if err != nil {
		global.LOG.Error("sjson set error", zap.Error(err))
		return errors.New("UnexpectedError")
	}
	// 将metadata上传到IPFS
	err, hash := IPFSUploadJSON(gjson.Parse(metadata).Value())
	if err != nil {
		global.LOG.Error("IPFSUploadJSON error", zap.Error(err))
		return errors.New("UnexpectedError")
	}
	// 构造请求
	data := DidCardRequest{
		Receiver: did,
		Params: DidCardParamsRequest{
			Title:       quest.Title,
			ChallengeID: cast.ToString(req.TokenId),
			Pass:        pass,
			Score:       score,
			Content:     "ipfs://" + hash,
		},
	}
	// 发送请求获取vc
	res, err := reqV3.C().R().SetBodyJsonMarshal(data).Post(global.CONFIG.ZCloak.Url + "/vc/issue")
	if err != nil {
		global.LOG.Error("get VC error", zap.Error(err))
		return errors.New("UnexpectedError")
	}
	if gjson.Get(res.String(), "code").Int() != 0 {
		global.LOG.Error("get VC error", zap.Error(err), zap.String("res", res.String()))
		return errors.New("UnexpectedError")
	}
	// 保存 VC
	zcloakCard := model.ZcloakCard{
		Address: address,
		Did:     did,
		QuestID: quest.ID,
		Score:   score,
		VC:      []byte(gjson.Get(res.String(), "data.vc").String()),
	}
	err = global.DB.Model(&model.ZcloakCard{}).Create(&zcloakCard).Error
	if err != nil {
		return err
	}
	err = SaveToNFTCollection(SaveCardInfoRequest{
		Chain:           "polygon",
		AccountAddress:  strings.ToLower(address),
		ContractAddress: strings.ToLower("0xc8e9cd4921e54c4163870092ca8d9660e967b53d"),
		TokenID:         cast.ToString(req.TokenId),
		ImageURI:        strings.TrimPrefix(gjson.Get(string(quest.MetaData), "image").String(), "ipfs://"),
		ErcType:         "erc1155",
		Name:            gjson.Get(string(quest.MetaData), "name").String(),
		DidAddress:      did,
	})
	if err != nil {
		return err
	}
	return
}

type SaveCardInfoRequest struct {
	Chain           string `json:"chain" form:"chain" binding:"required"`
	AccountAddress  string `json:"account_address" form:"account_address" binding:"required"`
	ContractAddress string `json:"contract_address" form:"contract_address" binding:"required"`
	TokenID         string `json:"token_id" form:"token_id" binding:"required"`
	ImageURI        string `json:"image_uri" form:"image_uri" binding:"required"`
	ErcType         string `json:"erc_type" form:"erc_type" binding:"required"`
	Name            string `json:"name" form:"name" binding:"required"`
	DidAddress      string `json:"did_address" form:"did_address" binding:"required"`
}

// SaveToNFTCollection 保存到NFT
func SaveToNFTCollection(saveCardInfo SaveCardInfoRequest) (err error) {
	if global.CONFIG.NFT.API == "" {
		return
	}
	// 发送请求
	client := reqV3.C().SetCommonHeader("x-api-key", global.CONFIG.NFT.APIKey)
	fmt.Println(global.CONFIG.NFT.API + "/zcloak/saveCardInfo")
	data, _ := json.Marshal(saveCardInfo)
	fmt.Println(string(data))
	r, err := client.R().SetBodyJsonMarshal(saveCardInfo).Post(global.CONFIG.NFT.API + "/zcloak/saveCardInfo")
	if err != nil {
		global.LOG.Error("SaveToNFT error", zap.Error(err), zap.String("res", r.String()))
		return err
	}
	if r.StatusCode != 200 {
		global.LOG.Error("SaveToNFT error", zap.Error(err), zap.String("res", r.String()))
		return errors.New("UnexpectedError")
	}
	if gjson.Get(r.String(), "status").Int() != 0 {
		global.LOG.Error("SaveToNFT error", zap.Error(err), zap.String("res", r.String()))
		return errors.New("UnexpectedError")
	}
	return nil
}
