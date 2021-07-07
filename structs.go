package mapstructurebenchmark

import "encoding/json"

func init() {
	employeeTypeRegistry["intern"] = NewIntern
	employeeTypeRegistry["salary"] = NewSalaryEmployee
}

type Employee interface {
	GetType() string
	GetName() string
	GetAge() int
}

type baseEmployee struct {
	Type string
	Name string
	Age  int
}

func (b baseEmployee) GetType() string {
	return b.Type
}

func (b baseEmployee) GetName() string {
	return b.Name
}

func (b baseEmployee) GetAge() int {
	return b.Age
}

type Intern struct {
	baseEmployee `mapstructure:",squash"`
	HourlyWage   float32
}

func NewIntern(data []byte) (Employee, error) {
	var i Intern
	err := json.Unmarshal(data, &i)
	return i, err
}

type SalaryEmployee struct {
	baseEmployee `mapstructure:",squash"`
	Salary       float32
}

func NewSalaryEmployee(data []byte) (Employee, error) {
	var s SalaryEmployee
	err := json.Unmarshal(data, &s)
	return s, err
}
