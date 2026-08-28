package controller

type CreateLinkDTO struct {
	DestinationURL string `json:"destination_url" validate:"required,url"`
	DisplayName    string `json:"display_name"    validate:"required"`
	Description    string `json:"description"`
}
