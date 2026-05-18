package dto

type CreateBookRequest struct {
	Title         string `json:"title" validate:"required,max=200"`
	Author        string `json:"author" validate:"required,max=100"`
	Description   string `json:"description" validate:"max=1000"`
	PublishedYear int    `json:"published_year" validate:"min=1000,max=9999"`
}

type UpdateBookRequest struct {
	Title         *string `json:"title" validate:"omitempty,max=200"`
	Author        *string `json:"author" validate:"omitempty,max=100"`
	Description   *string `json:"description" validate:"omitempty,max=1000"`
	PublishedYear *int    `json:"published_year" validate:"omitempty,min=1000,max=9999"`
}