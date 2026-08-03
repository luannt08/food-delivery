package domain

type Token struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
}

func NewToken(userId, refreshToken, createdAt string) *Token {
	return &Token{
		UserId:       userId,
		RefreshToken: refreshToken,
		CreatedAt:    createdAt,
	}
}

func (t *Token) VerifyRefreshToken(refreshToken string, userId string) bool {
	return t.RefreshToken == refreshToken && t.UserId == userId
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}