CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS employees (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	surname VARCHAR(50) NOT NULL,
    gender gender,
    age INTEGER CHECK (age >= 0),
	date_of_birth DATE,
    email VARCHAR(50),
    CONSTRAINT unique_name_surname_dateofbirth UNIQUE (name, surname, date_of_birth)
);

CREATE TABLE IF NOT EXISTS subscribers (
    id SERIAL PRIMARY KEY,
    employee_id INT REFERENCES employees(id),
    subscribed_to_id INT REFERENCES employees(id),
    UNIQUE (employee_id, subscribed_to_id)
);

