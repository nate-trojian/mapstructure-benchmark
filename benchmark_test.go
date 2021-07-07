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
	// Unmarshal into a map
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	// Get our type field
	t, ok := m["type"]
	if !ok {
		return nil, ErrEmployeeTypeUnknown
	}
	// Decode into the proper type
	switch t {
	case "intern":
		var i Intern
		err = mapstructure.Decode(m, &i)
		return i, err
	case "salary":
		var s SalaryEmployee
		err = mapstructure.Decode(m, &s)
		return s, err
	default:
		return nil, ErrEmployeeTypeUnknown
	}
}

func BenchmarkMapstructure(b *testing.B) {
	b.Run("Intern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := mapDecode(internInput); err != nil {
				b.Fatalf("Failed to convert intern - %v", err)
			}
		}
	})
	b.Run("Salary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := mapDecode(salariedInput); err != nil {
				b.Fatalf("Failed to convert salary - %v", err)
			}
		}
	})
}

func jsonConvertwRegistry(data []byte) (Employee, error) {
	// Unmarshal to base structure
	var base baseEmployee
	err := json.Unmarshal(data, &base)
	if err != nil {
		return nil, err
	}
	// Use base structure to get employee function from registry
	fn, err := GetEmployeeFunc(base.Type)
	if err != nil {
		return nil, err
	}

	// Convert
	return fn(data)
}

func BenchmarkJSONwRegistry(b *testing.B) {
	b.Run("Intern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := jsonConvertwRegistry(internInput); err != nil {
				b.Fatalf("Failed to convert intern - %v", err)
			}
		}
	})
	b.Run("Salary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := jsonConvertwRegistry(salariedInput); err != nil {
				b.Fatalf("Failed to convert salary - %v", err)
			}
		}
	})
}

func jsonConvertwSwitch(data []byte) (Employee, error) {
	// Unmarshal to base structure
	var base baseEmployee
	err := json.Unmarshal(data, &base)
	if err != nil {
		return nil, err
	}
	// Get the proper
	switch base.Type {
	case "intern":
		var i Intern
		err = json.Unmarshal(data, &i)
		return i, err
	case "salary":
		var s SalaryEmployee
		err = json.Unmarshal(data, &s)
		return s, err
	default:
		return nil, ErrEmployeeTypeUnknown
	}
}

func BenchmarkJSONwSwitch(b *testing.B) {
	b.Run("Intern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := jsonConvertwSwitch(internInput); err != nil {
				b.Fatalf("Failed to convert intern - %v", err)
			}
		}
	})
	b.Run("Salary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := jsonConvertwSwitch(salariedInput); err != nil {
				b.Fatalf("Failed to convert salary - %v", err)
			}
		}
	})
}
