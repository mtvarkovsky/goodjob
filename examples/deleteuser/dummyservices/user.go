package dummyservices

type (
	UserServiceClient struct {
		SafeDeleteUserSideEffect func() error
		RestoreUserSideEffect    func() error
	}

	SafeDeleteUserRequest struct {
		UserID string
	}

	RestoreUserRequest struct {
		UserID string
	}
)

func NewUserServiceClient(safeDeleteUserSideEffect func() error, restoreUserSideEffect func() error) *UserServiceClient {
	return &UserServiceClient{
		SafeDeleteUserSideEffect: safeDeleteUserSideEffect,
		RestoreUserSideEffect:    restoreUserSideEffect,
	}
}

func (s *UserServiceClient) SafeDeleteUser(request *SafeDeleteUserRequest) error {
	return s.SafeDeleteUserSideEffect()
}

func (s *UserServiceClient) RestoreUser(request *RestoreUserRequest) error {
	return s.RestoreUserSideEffect()
}
