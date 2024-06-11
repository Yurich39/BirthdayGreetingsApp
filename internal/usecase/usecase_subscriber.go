package usecase

import (
	"context"
	"fmt"

	"BirthdayGreetingsService/internal/entity"
)

type SubscriberUseCase struct {
	repo SubscribersRepo
	log  Logger
}

func NewSubscribers(repoActorsMovies SubscribersRepo, l Logger) *SubscriberUseCase {
	return &SubscriberUseCase{
		repo: repoActorsMovies,
		log:  l,
	}
}

func (uc *SubscriberUseCase) Add(ctx context.Context, data entity.Subscriber) error {

	err := uc.repo.Add(ctx, data)

	if err != nil {
		return fmt.Errorf("%s: repo.Add returned error: %w", op, err)
	}

	return nil
}

func (uc *SubscriberUseCase) Remove(ctx context.Context, employee_id, subscriber_id int) (entity.Employee, error) {

	res, err := uc.repo.Remove(ctx, employee_id, subscriber_id)

	if err != nil {
		return res, fmt.Errorf("%s: repo.Remove returned error: %w", op, err)
	}

	return res, nil
}

func (uc *SubscriberUseCase) ListAll(ctx context.Context, id int) ([]entity.Employee, error) {

	res, err := uc.repo.ListAll(ctx, id)

	if err != nil {
		return res, fmt.Errorf("%s: repo.List returned error: %w", op, err)
	}

	return res, nil
}

func (uc *SubscriberUseCase) ListHappyBirthdays(ctx context.Context, id int) ([]entity.Employee, error) {

	res, err := uc.repo.ListHappyBirthdays(ctx, id)

	if err != nil {
		return res, fmt.Errorf("%s: repo.List returned error: %w", op, err)
	}

	return res, nil
}
