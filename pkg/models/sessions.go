package models

type NewEmailSession struct {
	Email    string `json:"email" validate:"required" minLength:"5" maxLength:"255" example:"john.doe@email.com"`
	Password string `json:"password" minLength:"8"`
}

type NewDiscordSession struct {
	AuthCode string `json:"authCode" validate:"required"`
}
