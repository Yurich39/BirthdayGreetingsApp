package usecase

import (
	"context"

	"BirthdayGreetingsService/internal/entity"

	"golang.org/x/exp/slog"
)

type (
	Employee interface {
		Save(ctx context.Context, data entity.EmployeeData) (entity.Employee, error)
		Update(ctx context.Context, updates entity.Employee) (entity.Employee, error)
		Delete(ctx context.Context, id int) (entity.Employee, error)
		List(ctx context.Context) ([]entity.Employee, error)
		Next(ctx context.Context) ([]entity.Employee, error)
	}

	Subscriber interface {
		Add(ctx context.Context, data entity.Subscriber) error
		Remove(ctx context.Context, employee_id, subscriber_id int) (entity.Employee, error)
		ListAll(ctx context.Context, id int) ([]entity.Employee, error)
		ListHappyBirthdays(ctx context.Context, id int) ([]entity.Employee, error)
	}

	EmployeesRepo interface {
		Save(ctx context.Context, data entity.EmployeeData) (int, error)
		Update(ctx context.Context, updates entity.Employee) (entity.Employee, error)
		Delete(ctx context.Context, id int) (entity.Employee, error)
		List(ctx context.Context) ([]entity.Employee, error)
		Next(ctx context.Context) ([]entity.Employee, error)
	}

	SubscribersRepo interface {
		Add(ctx context.Context, data entity.Subscriber) error
		Remove(ctx context.Context, employee_id, subscriber_id int) (entity.Employee, error)
		ListAll(ctx context.Context, id int) ([]entity.Employee, error)
		ListHappyBirthdays(ctx context.Context, id int) ([]entity.Employee, error)
	}

	Logger interface {
		Debug(msg string, args ...any)
		Info(msg string, args ...any)
		Warn(msg string, args ...any)
		Error(msg string, args ...any)
		Err(err error) slog.Attr
	}
)
