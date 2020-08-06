package nextengine

import (
	"context"
)

type defaultTokenRepository struct {
	t Token
}

func (tr *defaultTokenRepository) Token(_ context.Context) (Token, error) {
	return tr.t, nil
}

func (tr *defaultTokenRepository) Save(_ context.Context, tok Token) error {
	token := tr.t

	if tok.AccessToken != "" {
		token.AccessToken = tok.AccessToken
	}
	if tok.RefreshToken != "" {
		token.RefreshToken = tok.RefreshToken
	}

	if tok.AccessTokenEndDate != "" {
		token.AccessTokenEndDate = tok.AccessTokenEndDate
	}
	if tok.RefreshTokenEndDate != "" {
		token.RefreshTokenEndDate = tok.RefreshTokenEndDate
	}

	return nil
}
