package repo

import (
	"context"
	"database/sql"
	"fmt"

	"BirthdayGreetingsService/internal/entity"

	"github.com/jmoiron/sqlx"
)

type SubscribersRepo struct {
	db *sqlx.DB
}

func NewSubscriberRepo(db *sql.DB) *SubscribersRepo {
	return &SubscribersRepo{db: sqlx.NewDb(db, "postgres")}
}

const SubscriberQueryAdd = `INSERT INTO subscribers (employee_id, subscribed_to_id) VALUES ($1, $2)`

func (r *SubscribersRepo) Add(ctx context.Context, data entity.Subscriber) error {

	//Проверяем наличие работника в таблице employees
	EmployeeQuery := `SELECT COUNT(*) FROM employees WHERE id = $1`

	var count1 int

	err := r.db.GetContext(ctx, &count1, EmployeeQuery,
		data.Employee_id,
	)

	if err != nil {
		return err
	}

	if count1 == 0 {
		return fmt.Errorf("employee_id was NOT found in database table 'employees'")
	}

	//Проверяем наличие работника в таблице employees
	SubscriberQuery := `SELECT COUNT(*) FROM employees WHERE id = $1`

	var count2 int

	err = r.db.GetContext(ctx, &count2, SubscriberQuery,
		data.Employee_id,
	)

	if err != nil {
		return err
	}

	if count2 == 0 {
		return fmt.Errorf("employee_id was NOT found in database table 'employees'")
	}

	// Вносим данные в базу данных в таблицу subscribers
	_, err = r.db.ExecContext(ctx, SubscriberQueryAdd,
		data.Employee_id,
		data.Subscribed_to_id,
	)

	if err != nil {
		return err
	}

	return nil
}

const ListEmployeesQuery = `SELECT 
	employees.id AS id,
	employees.name AS name,
    employees.surname AS surname,
    employees.gender AS gender,
    employees.age AS age,
	employees.email AS email,
    TO_CHAR(employees.date_of_birth, 'DD.MM.YYYY') AS date_of_birth
	FROM 
    employees
	JOIN 
    subscribers ON subscribers.subscribed_to_id = employees.id
	WHERE 
    subscribers.employee_id = $1`

func (r *SubscribersRepo) ListAll(ctx context.Context, id int) ([]entity.Employee, error) {
	var data []entity.Employee
	err := r.db.SelectContext(ctx, &data, ListEmployeesQuery, id)

	if err != nil {
		return nil, fmt.Errorf("%s: DB returned error: %w", op, err)
	}

	return data, nil
}

const ListEmployeesWithHappyBirthdaysQuery = `SELECT 
    employees.id,
    employees.name,
    employees.surname,
    employees.gender,
    employees.age,
	employees.email AS email,
    TO_CHAR(employees.date_of_birth, 'DD.MM.YYYY') AS date_of_birth
	FROM 
    employees
	JOIN 
    subscribers ON subscribers.subscribed_to_id = employees.id
	WHERE 
    subscribers.employee_id = $1 
    AND EXTRACT(MONTH FROM employees.date_of_birth) = EXTRACT(MONTH FROM CURRENT_DATE)
    AND EXTRACT(DAY FROM employees.date_of_birth) = EXTRACT(DAY FROM CURRENT_DATE)`

func (r *SubscribersRepo) ListHappyBirthdays(ctx context.Context, id int) ([]entity.Employee, error) {
	var data []entity.Employee
	err := r.db.SelectContext(ctx, &data, ListEmployeesWithHappyBirthdaysQuery, id)

	if err != nil {
		return nil, fmt.Errorf("%s: DB returned error: %w", op, err)
	}

	return data, nil
}

const DeleteSubscriber = `DELETE FROM subscribers
	WHERE employee_id = $1
	AND subscribed_to_id = $2`

func (r *SubscribersRepo) Remove(ctx context.Context, employee_id, subscriber_id int) (entity.Employee, error) {

	_, err := r.db.ExecContext(ctx, DeleteSubscriber, employee_id, subscriber_id)

	if err != nil {
		return entity.Employee{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return entity.Employee{Id: &subscriber_id}, nil
}
