package auth

type RegisterInput struct {
	Username string `binding:"required"       json:"username"`
	Password string `binding:"required"       json:"password"`
	Email    string `binding:"required,email" json:"email"`
}

type LoginInput struct {
	Username string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
}

type OauthInput struct {
	Code string `binding:"required" json:"code"`
}
