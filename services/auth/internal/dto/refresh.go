package dto

type RefreshInput struct {
    // string refresh_token = 1;
	RefreshToken string
}

type RefreshOutput struct {
    // string access_token = 1;
    // string refresh_token = 2;
	AccessToken  string
	RefreshToken string
}