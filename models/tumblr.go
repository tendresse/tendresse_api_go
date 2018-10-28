package models

type Tumblr struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Response struct {
		Blog struct {
			Title      string `json:"title"`
			TotalPosts int    `json:"total_posts"`
		} `json:"blog"`
		Posts []struct {
			Id     int      `json:"id"`
			Type   string   `json:"type"`
			Tags   []string `json:"tags"`
			Photos []struct {
				OriginalSize struct {
					Url string `json:"url"`
				} `json:"original_size"`
			} `json:"photos"`
		} `json:"posts"`
	} `json:"response"`
}
