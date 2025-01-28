package request

type CreateLink struct {
	SrcLink string `json:"link" binding:"required"`

}