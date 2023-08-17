package backend

import (
	"backend/internal/app/model/receive"
	"backend/internal/app/model/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/url"
	"time"
)

// GetYouTubePlayList 获取 YouTube 播放列表
func GetYouTubePlayList(link string) (result []response.GetYouTubePlayListResponse, err error) {
	baseUrl := "https://www.googleapis.com/youtube/v3/playlistItems"
	key := "AIzaSyARUWxJ6NuHic1g1HnRKK8UAxib-2W3mKo"
	parsedURL, err := url.Parse(link)
	if err != nil {
		fmt.Printf("URL 解析错误：%s\n", err)
		return
	}
	queryParams := parsedURL.Query()
	listID := queryParams.Get("list")
	if listID == "" {
		return result, errors.New("请输入正确的链接")
	}
	params := map[string]string{"part": "snippet", "playlistId": listID, "key": key, "maxResults": "50"}
	client := req.C().SetTimeout(30 * time.Second)
	res, err := client.R().SetQueryParams(params).Get(baseUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	var playlistItems receive.PlaylistItems
	err = json.Unmarshal([]byte(res.String()), &playlistItems)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(playlistItems.Items) == 0 {
		return result, errors.New("解析失败，请检查链接是否正确")
	}
	for _, v := range playlistItems.Items {
		url := fmt.Sprintf("https://www.youtube.com/watch?v=%s&list=%s", v.Snippet.ResourceID.VideoID, listID)
		content := response.GetYouTubePlayListResponse{
			Label: v.Snippet.Title,
			ID:    v.Snippet.ResourceID.VideoID,
			Img:   v.Snippet.Thumbnails.Maxres.URL,
			Url:   url,
		}
		result = append(result, content)
	}
	return result, nil
}
