package impl

import (
	"errors"
	"github.com/tangvis/erp/agent/mysql"
	"testing"
)

// Define a mock status type to match define.Status in your real struct
type MockStatus int

const (
	StatusActive MockStatus = iota
	StatusInactive
)

// Mock the BrandTab with S interface implemented for testing
type MockBrandTab struct {
	Name        string     `al:"品牌名"`
	Desc        string     `al:"品牌描述"`
	URL         string     // URLTab `al:"图片链接"`
	BrandStatus MockStatus `al:"状态"`
	CreateBy    string     // email of creator

	mysql.BaseModel
}

type OtherTab struct {
}

func (m MockStatus) String() string {
	switch m {
	case StatusActive:
		return "active"
	default:
		return "inactive"
	}
}

func TestCompare(t *testing.T) {
	testCases := []struct {
		name     string
		before   any
		after    any
		expected map[string]string
		err      error
	}{
		{
			name:     "identical structs",
			before:   MockBrandTab{Name: "Brand A", Desc: "Description A", BrandStatus: StatusActive},
			after:    MockBrandTab{Name: "Brand A", Desc: "Description A", BrandStatus: StatusActive},
			expected: map[string]string{},
			err:      nil,
		},
		{
			name:   "different structs",
			before: MockBrandTab{Name: "Brand A", Desc: "Description A", BrandStatus: StatusActive},
			after:  MockBrandTab{Name: "Brand B", Desc: "Description B", BrandStatus: StatusInactive},
			expected: map[string]string{
				"品牌名":  "[Brand A] has been changed to [Brand B]",
				"品牌描述": "[Description A] has been changed to [Description B]",
				"状态":   "[active] has been changed to [inactive]",
			},
			err: nil,
		},
		{
			name:     "nil before struct",
			before:   nil,
			after:    MockBrandTab{Name: "Brand B"},
			expected: nil,
			err:      errors.New("one of the input values is nil"),
		},
		{
			name:     "nil after struct",
			before:   MockBrandTab{Name: "Brand A"},
			after:    nil,
			expected: nil,
			err:      errors.New("one of the input values is nil"),
		},
		{
			name:     "type mismatch",
			before:   MockBrandTab{Name: "Brand A"},
			after:    OtherTab{},
			expected: nil,
			err:      errors.New("type mismatch between before and after"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Compare(tc.before, tc.after)
			if (err != nil && tc.err == nil) || (err == nil && tc.err != nil) || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if len(result) != len(tc.expected) {
				t.Errorf("Expected result length %d, got %d", len(tc.expected), len(result))
			}
			for key, val := range tc.expected {
				if result[key] != val {
					t.Errorf("Expected %s to be '%s', got '%s'", key, val, result[key])
				}
			}
		})
	}
}
