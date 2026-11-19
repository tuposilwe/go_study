package service

import (
	"banking/domain"
	"banking/dto"
	"banking/errs"
)

type CustomerService interface {
	GetAllCustomers(status string) ([]dto.CustomerResponse, error)
	GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repo: repository}
}

func (s DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, error) {
	customers, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		responses = append(responses, c.ToDto())
	}

	return responses, nil
}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}

	response := c.ToDto()
	return &response, nil
}