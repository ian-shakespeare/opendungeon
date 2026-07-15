package models

type Provider struct {
	Name    string `json:"name" validate:"required" example:"Discord"`
	AuthURL string `json:"authUrl" validate:"required"`
}

type ProviderList struct {
	State     string     `json:"state" validate:"required"`
	Providers []Provider `json:"providers" validate:"required"`
}
