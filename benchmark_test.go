package mapstructurebenchmark

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
)

var (
	internMap = map[string]interface{}{
		"type":        "intern",
		"name":        "Intern 1",
		"age":         21,
		"hourly_wage": 20.0,
	}
	internInput, _ = json.Marshal(internMap)
	salariedMap    = map[string]interface{}{
		"type":   "salary",
		"name":   "Alice",
		"age":    30,
		"salary": 100000.0,
	}
	salariedInput, _ = json.Marshal(salariedMap)
)

func mapDecode(data []byte) (Employee, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	t, ok := m["type"]
	if !ok {
		return nil, ErrEmployeeTypeUnknown
	}
	switch t {
	case "intern":
		var i Intern
		err = mapstructure.Decode(m, &i)
		return i, err
	case "salary":
		var s SalaryEmployee
		err = mapstructure.Decode(m, &s)
		return s, err
	}
	return nil, ErrEmployeeTypeUnknown
}

func BenchmarkMapstructure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := mapDecode(internInput); err != nil {
			b.Fatalf("Failed to convert intern - %v", err)
		}
		if _, err := mapDecode(salariedInput); err != nil {
			b.Fatalf("Failed to convert salary - %v", err)
		}
	}
}

func jsonConvert(data []byte) (Employee, error) {
	var base baseEmployee
	err := json.Unmarshal(data, &base)
	if err != nil {
		return nil, err
	}
	fn, err := GetEmployeeFunc(base.Type)
	if err != nil {
		return nil, err
	}

	return fn(data)
}

func BenchmarkJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := jsonConvert(internInput); err != nil {
			b.Fatalf("Failed to convert intern - %v", err)
		}
		if _, err := jsonConvert(salariedInput); err != nil {
			b.Fatalf("Failed to convert salary - %v", err)
		}
	}
}
