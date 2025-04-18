package dto

type CreateLinkDTO struct {
	URL   string   `json:"url" binding:"required,url"`
	Title string   `json:"title" binding:"required"`
	Tags  []string `json:"tags"`
}

type UpdateLinkDTO struct {
	URL   *string   `json:"url" binding:"omitempty,url"` // optionnel mais validé s’il est là
	Title *string   `json:"title" binding:"omitempty"`   // idem
	Tags  *[]string `json:"tags"`                        // facultatif
}
