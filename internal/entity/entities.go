package entity

type Employee struct {
	Id *int `db:"id" json:"id,omitempty"`
	EmployeeData
}

type EmployeeData struct {
	Name        *string `db:"name" json:"name,omitempty"`
	Surname     *string `db:"surname" json:"surname,omitempty"`
	Gender      *string `db:"gender" json:"gender,omitempty"`
	Age         *int    `db:"age" json:"age,omitempty"`
	Email       *string `db:"email" json:"email,omitempty"`
	DateOfBirth *string `db:"date_of_birth" json:"date_of_birth,omitempty"`
}

type Subscriber struct {
	Id               *int `db:"id" json:"id,omitempty"`
	Employee_id      *int `db:"employee_id" json:"employee_id,omitempty"`
	Subscribed_to_id *int `db:"subscribed_to_id" json:"subscribed_to_id,omitempty"`
}
