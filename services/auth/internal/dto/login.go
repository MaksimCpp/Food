package dto

type LoginInput struct {
	Email    string
	Password string
}

type LoginAdminInput struct {
	Email     string
	Password  string
	AdminCode string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}