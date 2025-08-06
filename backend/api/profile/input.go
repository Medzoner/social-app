package profile

type UpdateProfileInput struct {
	Bio    *string `json:"bio,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
	Email  *string `json:"email,omitempty"`
}

type VerifyCodeInput struct {
	Code   string `binding:"required" json:"code"`
	Type   string `binding:"required" json:"type"`
	UserID uint64 `binding:"required" json:"id"`
}

type ResendCodeInput struct {
	Type   string `binding:"required" json:"type"`
	UserID uint64 `binding:"required" json:"id"`
}
