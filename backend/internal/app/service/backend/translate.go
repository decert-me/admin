package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/utils"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"path/filepath"
	"strings"
)

// SubmitTranslate GitHub Action 提交翻译
func SubmitTranslate(req request.SubmitTranslateRequest) (err error) {
	filePath := req.Filename
	// 从文件名获取TokenID
	filenameWithExt := filepath.Base(filePath)           // 获取文件名，包含后缀
	ext := filepath.Ext(filePath)                        // 获取文件后缀
	filename := strings.TrimSuffix(filenameWithExt, ext) // 去掉后缀
	// 获取Github英文翻译文件
	contentEn, err := getGithubTranslateFile("/translation/en/" + filenameWithExt)
	if err != nil {
		return errors.New("获取Github英文翻译文件失败")
	}
	// 中文翻译文件
	contentCn, err := getGithubTranslateFile("/translation/cn/" + filenameWithExt)
	if err != nil {
		return errors.New("获取Github中文翻译文件失败")
	}
	// 合辑处理
	if strings.Contains(filename, "collection") {
		id, err := cast.ToInt64E(strings.Replace(filename, "collection_", "", 1))
		if err != nil {
			return err
		}
		// 处理翻译文件
		collectionTranslate, err := handleTranslateCollection(id, contentEn)
		if err != nil {
			return err
		}
		// 保存翻译结果
		collectionTranslate.CollectionID = id
		collectionTranslate.Language = "en-US"
		err = saveTranslateResultCollection(id, collectionTranslate)
		if err != nil {
			return err
		}
		// 处理翻译文件
		collectionTranslate, err = handleTranslateCollection(id, contentCn)
		if err != nil {
			return err
		}
		// 保存翻译结果
		collectionTranslate.CollectionID = id
		collectionTranslate.Language = "zh-CN"
		err = saveTranslateResultCollection(id, collectionTranslate)
		return err
	}
	// 挑战处理
	tokenID, err := cast.ToInt64E(filename)
	if err != nil {
		return err
	}
	// 处理翻译文件
	questTranslate, err := handleTranslate(tokenID, contentEn)
	if err != nil {
		return
	}
	// 保存翻译结果
	questTranslate.TokenId = tokenID
	questTranslate.Language = "en-US"
	err = saveTranslateResult(tokenID, questTranslate)
	if err != nil {
		return err
	}

	// 处理翻译文件
	questTranslate, err = handleTranslate(tokenID, contentCn)
	if err != nil {
		return
	}
	// 保存翻译结果
	questTranslate.TokenId = tokenID
	questTranslate.Language = "zh-CN"
	err = saveTranslateResult(tokenID, questTranslate)
	return err
}

// 处理翻译文件
func handleTranslate(tokenID int64, content string) (questTranslate model.QuestTranslated, err error) {
	// 获取数据库Quest数据
	db := global.DB
	var quest model.Quest
	err = db.Where("token_id = ?", tokenID).First(&quest).Error
	if err != nil {
		return questTranslate, err
	}
	// 获取Title和Description
	questTranslate.Title = gjson.Get(content, "title").String()
	questTranslate.Description = gjson.Get(content, "description").String()
	// 处理翻译内容
	questRes, answerRes, err := handleTranslateContent(string(quest.QuestData), content)
	if err != nil {
		return questTranslate, err
	}
	// 处理MetaData
	metaDataRes, err := handleTranslateMetaData(string(quest.MetaData), content)
	if err != nil {
		return questTranslate, err
	}
	questTranslate.QuestData = []byte(questRes)
	questTranslate.MetaData = []byte(metaDataRes)
	questTranslate.Answer = answerRes
	return questTranslate, nil
}

// 处理翻译文件，合辑
func handleTranslateCollection(id int64, content string) (collectionTranslate model.CollectionTranslated, err error) {
	// 获取数据库Collection数据
	db := global.DB
	var collection model.Collection
	err = db.Where("id = ?", id).First(&collection).Error
	if err != nil {
		return collectionTranslate, err
	}
	// 获取Title和Description
	collectionTranslate.Title = gjson.Get(content, "title").String()
	collectionTranslate.Description = gjson.Get(content, "description").String()
	return collectionTranslate, nil
}

// 处理翻译内容
func handleTranslateContent(contentEN string, contentTr string) (content string, answerRes string, err error) {
	content = contentEN
	// 处理翻译内容
	// title
	content, err = sjson.Set(content, "title", gjson.Get(contentTr, "title").String())
	if err != nil {
		return "", "", err
	}
	// description
	content, err = sjson.Set(content, "description", gjson.Get(contentTr, "description").String())
	if err != nil {
		return "", "", err
	}
	// questions
	questions := gjson.Get(contentTr, "questions").Array()
	// 获取 answer
	answerRes = utils.AnswerDecode(global.CONFIG.Quest.EncryptKey, gjson.Get(contentEN, "answers").String())
	// 解密答案
	for index, question := range questions {
		//fmt.Println("index", index)
		//fmt.Println("question", question)
		// title
		titleContent := gjson.Get(question.String(), "title").String()
		content, err = sjson.Set(content, "questions."+cast.ToString(index)+".title", titleContent)
		if err != nil {
			return "", "", err
		}
		// options
		optionsContent := gjson.Get(question.String(), "options")
		if optionsContent.Exists() {
			//fmt.Println("optionsContent", optionsContent)
			content, err = sjson.SetRaw(content, "questions."+cast.ToString(index)+".options", optionsContent.Raw)
			// 填空题特殊处理
			if gjson.Get(question.String(), "type").String() == "fill_blank" {
				answerRes, _ = sjson.Set(answerRes, cast.ToString(index), optionsContent.Array()[0].String())
			}
		}
		// description
		description := gjson.Get(question.String(), "description")
		if description.Exists() {
			content, err = sjson.Set(content, "questions."+cast.ToString(index)+".description", description.String())
			if err != nil {
				return "", "", err
			}
		}
	}
	// 加密答案
	answerRes = utils.AnswerEncode(global.CONFIG.Quest.EncryptKey, answerRes)
	return content, answerRes, nil
}

// handleTranslateMetaData
func handleTranslateMetaData(metaDataEN string, contentTr string) (metaData string, err error) {
	metaData = metaDataEN
	// 处理翻译内容
	// title
	metaData, err = sjson.Set(metaData, "name", gjson.Get(contentTr, "title").String())
	if err != nil {
		return "", err
	}
	// description
	metaData, err = sjson.Set(metaData, "description", gjson.Get(contentTr, "description").String())
	if err != nil {
		return "", err
	}
	// challenge_title
	metaData, err = sjson.Set(metaData, "attributes.challenge_title", gjson.Get(contentTr, "title").String())
	return metaData, nil
}

// getGithubTranslateFile 获取Github翻译文件
func getGithubTranslateFile(filePath string) (content string, err error) {
	client := req.C()
	repo := strings.Replace(global.CONFIG.Translate.GithubRepo, "https://github.com", "https://raw.githubusercontent.com", 1)
	baseURL := repo + "/" + global.CONFIG.Translate.GithubBranch
	fmt.Println(baseURL + filePath)
	res, err := client.R().Get(baseURL + filePath)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("获取翻译文件失败")
	}
	return res.String(), nil
}

// saveTranslateResult 保存翻译结果
func saveTranslateResult(tokenID int64, questTranslate model.QuestTranslated) (err error) {
	db := global.DB
	// 已存在则更新
	var count int64
	err = db.Model(&model.QuestTranslated{}).Where("token_id = ? AND language = ?", tokenID, questTranslate.Language).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		err = db.Model(&model.QuestTranslated{}).Where("token_id = ? AND language = ?", tokenID, questTranslate.Language).Updates(&questTranslate).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		// 不存在则创建
		err = db.Model(&model.QuestTranslated{}).Create(&questTranslate).Error
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// saveTranslateResultCollection 保存翻译结果
func saveTranslateResultCollection(id int64, collectionTranslate model.CollectionTranslated) (err error) {
	db := global.DB
	// 已存在则更新
	var count int64
	err = db.Model(&model.CollectionTranslated{}).Where("id = ? AND language = ?", id, collectionTranslate.Language).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		err = db.Model(&model.CollectionTranslated{}).Where("id = ? AND language = ?", id, collectionTranslate.Language).Updates(&collectionTranslate).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		// 不存在则创建
		err = db.Model(&model.CollectionTranslated{}).Create(&collectionTranslate).Error
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
