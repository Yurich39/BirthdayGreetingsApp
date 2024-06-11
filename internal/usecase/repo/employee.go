package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"BirthdayGreetingsService/internal/controller/middleware/filter"
	"BirthdayGreetingsService/internal/controller/middleware/pagination"
	"BirthdayGreetingsService/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const op = "internal.usecase.repo"
const pageSize uint64 = 10

type EmplyeesRepo struct {
	db *sqlx.DB
}

func NewEmployeesRepo(db *sql.DB) *EmplyeesRepo {
	return &EmplyeesRepo{db: sqlx.NewDb(db, "postgres")}
}

const EmployeeQuerySave = `INSERT INTO employees(name, surname, gender, age, email, date_of_birth)
					VALUES($1, $2, $3, $4, $5, TO_DATE($6, 'DD.MM.YYYY'))
					ON CONFLICT (name, surname, date_of_birth) DO NOTHING
    				RETURNING id`

func (r *EmplyeesRepo) Save(ctx context.Context, data entity.EmployeeData) (int, error) {

	var res int

	err := r.db.GetContext(ctx, &res, EmployeeQuerySave,
		data.Name,
		data.Surname,
		data.Gender,
		data.Age,
		data.Email,
		data.DateOfBirth,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *EmplyeesRepo) Update(ctx context.Context, updates entity.Employee) (entity.Employee, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Update("employees")

	// Составим выражение для оператора SQL SET
	data, err := getMapEmployee(updates)

	if err != nil {
		return entity.Employee{}, fmt.Errorf("%s: Error: %w", op, err)
	}

	qb = qb.SetMap(data)

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", *updates.Id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, name, surname, gender, age, email, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth"

	if err != nil {
		return entity.Employee{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var res entity.Employee

	err = r.db.GetContext(ctx, &res, sql, i...)

	if err != nil {
		return entity.Employee{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *EmplyeesRepo) Delete(ctx context.Context, id int) (entity.Employee, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Delete("employees")

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, name, surname, gender, age, email, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth"

	if err != nil {
		return entity.Employee{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var res entity.Employee

	err = r.db.GetContext(ctx, &res, sql, i...)

	if err != nil {
		return entity.Employee{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *EmplyeesRepo) List(ctx context.Context) ([]entity.Employee, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, gender, age, email, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth").From("employees")

	// Составим выражение для оператора SQL Where ... AND ...
	filter_options, _ := ctx.Value(filter.FilterOptionsContextKey).(map[string][]string)

	stmt := []string{}
	for k, v := range filter_options {
		for _, val := range v {
			stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, val))
		}
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	// Используем оператор ORDER BY
	qb = qb.OrderBy("id ASC")

	qb = qb.Limit(pageSize)

	sql, i, err := qb.ToSql()

	if err != nil {
		return []entity.Employee{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Employee{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Employee{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *EmplyeesRepo) Next(ctx context.Context) ([]entity.Employee, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, gender, age, email, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth").From("employees")

	// Составим выражение для оператора SQL Where ... AND ...
	filter_options, _ := ctx.Value(filter.FilterOptionsContextKey).(map[string][]string)

	stmt := []string{}
	for k, v := range filter_options {
		for _, val := range v {
			stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, val))
		}
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	// Условие пагинации
	personID := ctx.Value(pagination.NextPersonID).(int)
	qb = qb.Where(fmt.Sprintf("id >= %d", personID))

	// Используем оператор ORDER BY
	qb = qb.OrderBy("id ASC")

	qb = qb.Limit(pageSize)

	sql, i, err := qb.ToSql()

	if err != nil {
		return []entity.Employee{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Employee{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Employee{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func getMapEmployee(updates entity.Employee) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	if val := updates.EmployeeData.Name; val != nil {
		res["name"] = *val
	}

	if val := updates.EmployeeData.Surname; val != nil {
		res["surname"] = *val
	}

	if val := updates.EmployeeData.Gender; val != nil {
		res["gender"] = *val
	}

	if val := updates.EmployeeData.Age; val != nil {
		res["age"] = *val
	}

	if val := updates.EmployeeData.Email; val != nil {
		res["email"] = *val
	}

	if val := updates.EmployeeData.DateOfBirth; val != nil {
		layout := "02.01.2006" // шаблон формата даты в строке

		// распарсим строку в формат даты
		date, err := time.Parse(layout, *val)
		if err != nil {
			return res, fmt.Errorf("%s: Ошибка при парсинге даты", err)
		}

		// Преобразовать дату в формат PostgreSQL
		res["date_of_birth"] = date.Format("2006-01-02")
	}

	if len(res) == 0 {
		return res, fmt.Errorf("%s: Data for update operation were NOT specified", op)
	}

	return res, nil
}
