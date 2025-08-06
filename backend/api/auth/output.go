package auth

type JWTToken struct {
	AccessToken    string         `json:"access_token"`
	RefreshToken   string         `json:"refresh_token"`
	IDToken        string         `json:"id_token"`
	VerifyResponse VerifyResponse `json:"-"`
}

type VerifyResponse struct {
	Username string `json:"username,omitempty"`
	ID       uint64 `json:"id,omitempty"`
	Verified bool   `json:"verified,omitempty"`
}

func (v VerifyResponse) IsEmpty() bool {
	return v.ID == 0 && v.Username == "" && !v.Verified
}
