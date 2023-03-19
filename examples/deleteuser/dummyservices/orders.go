package dummyservices

type (
	OrdersServiceClient struct {
		SafeDeleteUserOrdersSideEffect func() error
		RestoreUserOrdersSideEffect    func() error
	}

	SafeDeleteUserOrdersRequest struct {
		UserID string
	}

	RestoreUserOrdersRequest struct {
		UserID string
	}
)

func NewOrdersServiceClient(safeDeleteUserOrdersSideEffect func() error, restoreUserOrdersSideEffect func() error) *OrdersServiceClient {
	return &OrdersServiceClient{
		SafeDeleteUserOrdersSideEffect: safeDeleteUserOrdersSideEffect,
		RestoreUserOrdersSideEffect:    restoreUserOrdersSideEffect,
	}
}

func (s *OrdersServiceClient) SafeDeleteUserOrders(request *SafeDeleteUserOrdersRequest) error {
	return s.SafeDeleteUserOrdersSideEffect()
}

func (s *OrdersServiceClient) RestoreUserOrders(request *RestoreUserOrdersRequest) error {
	return s.RestoreUserOrdersSideEffect()
}
