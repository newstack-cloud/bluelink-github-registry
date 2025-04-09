package auth

import "context"

// TokenService is an interface that defines the methods
// required for an authentication token service.
type TokenService interface {
	// ValidateToken validates the provided token and returns
	// true if the token is valid, false otherwise.
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type githubTokenService struct {
}

func NewGitHubTokenService() TokenService {
	return &githubTokenService{}
}

func (s *githubTokenService) ValidateToken(ctx context.Context, token string) (bool, error) {
	return false, nil
}
