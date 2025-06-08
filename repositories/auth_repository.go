package repositories

type AuthRepository struct{}

func (r AuthRepository) CheckCustomerExists(customerID string) bool {
	return customerID != ""
}