package response

type GetLink struct {
	SrcLink   string `json:"src_link"`
	ShortLink string `json:"short_link"`
	Visits    int    `json:"count_visits"`
}
