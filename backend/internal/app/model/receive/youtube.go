package receive

import "time"

type PlaylistItems struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []struct {
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Maxres struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
			PlaylistID   string `json:"playlistId"`
			Position     int    `json:"position"`
			ResourceID   struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
			VideoOwnerChannelTitle string `json:"videoOwnerChannelTitle"`
			VideoOwnerChannelID    string `json:"videoOwnerChannelId"`
		} `json:"snippet"`
	} `json:"items"`
}
