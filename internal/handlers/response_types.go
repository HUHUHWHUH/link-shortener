package handlers

type ShorUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type OriginalUrlResponse struct {
	Url string `json:"url"`
}
