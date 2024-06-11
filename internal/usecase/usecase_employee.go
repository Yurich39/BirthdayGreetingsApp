package usecase

import (
	"context"
	"fmt"

	"BirthdayGreetingsService/internal/entity"
)

const op = "internal.usecase"

type EmployeeUseCase struct {
	repo EmployeesRepo
	log  Logger
}

func NewEmployees(repoEmployees EmployeesRepo, l Logger) *EmployeeUseCase {
	return &EmployeeUseCase{
		repo: repoEmployees,
		log:  l,
	}
}

func (uc *EmployeeUseCase) Save(ctx context.Context, data entity.EmployeeData) (entity.Employee, error) {

	id, err := uc.repo.Save(ctx, data)

	if err != nil {
		return entity.Employee{}, err
	}

	res := entity.Employee{
		Id:        &id,
		EmployeeData: data,
	}

	return res, nil
}

func (uc *EmployeeUseCase) Update(ctx context.Context, updates entity.Employee) (entity.Employee, error) {
	res, err := uc.repo.Update(ctx, updates)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Update returned error: %w", op, err)
	}

	return res, nil
}

func (uc *EmployeeUseCase) Delete(ctx context.Context, id int) (entity.Employee, error) {
	res, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Delete returned error: %w", op, err)
	}

	return res, nil
}

func (uc *EmployeeUseCase) List(ctx context.Context) ([]entity.Employee, error) {
	res, err := uc.repo.List(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.List returned error: %w", op, err)
	}

	return res, nil
}

func (uc *EmployeeUseCase) Next(ctx context.Context) ([]entity.Employee, error) {
	res, err := uc.repo.Next(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Next returned error: %w", op, err)
	}

	return res, nil
}
