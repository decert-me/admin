package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/utils"
	"encoding/json"
	"errors"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"strconv"
)

func RunAirdrop(r request.RunAirdropReq) (err error) {
	// 生成params
	paramsMap := map[string]interface{}{
		"app": r.App,
	}
	// 将Map转换为JSON格式的字节数组
	paramsData, err := json.Marshal(paramsMap)
	if err != nil {
		return
	}
	var body struct {
		Params string `json:"params"`
	}
	body.Params = string(paramsData)

	url := global.CONFIG.Airdrop.Api + "/v1/airdrop/runAirdrop"
	// 生成校验hash和时间戳
	timestamp, hashValue := utils.HashData(body, global.CONFIG.Airdrop.VerifyKey)
	headers := map[string]string{
		"verify":    hashValue,
		"timestamp": strconv.Itoa(int(timestamp)),
	}
	client := req.C()
	res, err := client.R().SetHeaders(headers).SetBodyJsonMarshal(body).Post(url)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		return errors.New("error")
	}
	if gjson.Get(res.String(), "status").Int() != 0 {
		return errors.New(gjson.Get(res.String(), "message").String())
	}
	return nil
}
