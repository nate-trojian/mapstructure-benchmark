package mapstructurebenchmark

import "errors"

type EmployeeFunc func([]byte) (Employee, error)

var (
	employeeTypeRegistry   = make(map[string]EmployeeFunc)
	ErrEmployeeTypeUnknown = errors.New("employee type unknown")
)

func GetEmployeeFunc(t string) (EmployeeFunc, error) {
	fn, ok := employeeTypeRegistry[t]
	if !ok {
		return nil, ErrEmployeeTypeUnknown
	}
	return fn, nil
}
