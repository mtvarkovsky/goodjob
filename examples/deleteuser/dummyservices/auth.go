package dummyservices

type (
	AuthServiceClient struct {
		SafeDeleteAuthDataSideEffect func() error
		RestoreAuthDataSideEffect    func() error
	}

	SafeDeleteAuthDataRequest struct {
		UserID string
	}

	RestoreAuthDataRequest struct {
		UserID string
	}
)

func NewAuthServiceClient(safeDeleteAuthDataSideEffect func() error, restoreAuthDataSideEffect func() error) *AuthServiceClient {
	return &AuthServiceClient{
		SafeDeleteAuthDataSideEffect: safeDeleteAuthDataSideEffect,
		RestoreAuthDataSideEffect:    restoreAuthDataSideEffect,
	}
}

func (s *AuthServiceClient) SafeDeleteAuthData(request *SafeDeleteAuthDataRequest) error {
	return s.SafeDeleteAuthDataSideEffect()
}

func (s *AuthServiceClient) RestoreAuthData(request *RestoreAuthDataRequest) error {
	return s.RestoreAuthDataSideEffect()
}
