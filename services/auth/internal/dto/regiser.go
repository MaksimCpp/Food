package dto

type RegisterInput struct {
    // string username = 1;
    // string email = 2;
    // string password = 3;
	Username string
	Email    string
	Password string
}

type RegisterAdminInput struct {
    // string username = 1;
    // string email = 2;
    // string password = 3;
    // string admin_code = 4;
	Username  string
	Email     string
	Password  string
	AdminCode string
}

type RegisterOutput struct {
    // int64 user_id = 1;
    // string email = 2;
	UserID int64
	Email  string
}