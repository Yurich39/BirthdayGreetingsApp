package api

import "BirthdayGreetingsService/internal/entity"

type EmployeeResponse struct {
	Status         string            `json:"status,omitempty"`
	Employee       *entity.Employee  `json:"employee,omitempty"`
	Employees      []entity.Employee `json:"employees,omitempty"`
	NextEmployeeID int               `json:"next_employee_id,omitempty"`
}

type SubscriberResponse struct {
	Status     string             `json:"status,omitempty"`
	Subscriber *entity.Subscriber `json:"subscriber,omitempty"`
}

type SubscribersResponse struct {
	Status    string            `json:"status,omitempty"`
	Employees []entity.Employee `json:"employees,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusEmpty = "No subscribers with birthdays today"
	StatusError = "Error"
)

type customError struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Error(msg string) customError {
	return customError{
		Status: StatusError,
		Error:  msg,
	}
}
