package dto

type CreateLinkDTO struct {
	URL   string   `json:"url" binding:"required,url"`
	Title string   `json:"title" binding:"required"`
	Tags  []string `json:"tags"`
}
