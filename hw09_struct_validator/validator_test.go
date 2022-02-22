package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

type TestError struct {
	Field string
	Err   error
}

func (e TestError) Error() string {
	return e.Field + ": " + e.Err.Error() + "\n"
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				"111",
				"xxx", 33, "test@test.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			TestError{
				"ID",
				errors.New(
					"the length of 3 of the value '111' is not equal to 36: not equal: validation error",
				),
			},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 3, "test@test.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			TestError{
				"Age",
				errors.New(
					"3 less then 18: min value exceeded: validation error",
				),
			},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 333, "test@test.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			TestError{
				"Age",
				errors.New(
					"333 more then 50: max value exceeded: validation error",
				),
			},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "testtest.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			TestError{
				"Email",
				errors.New(
					"the 'testtest.ru' value does not match the '^\\w+@\\w+\\.\\w+$' pattern: not match: validation error",
				),
			},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "test@test.ru", "not",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			TestError{
				"Role",
				errors.New(
					"the value 'not' is not included in the list of 'admin,stuff': not contained: validation error",
				),
			},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "test@test.ru", "stuff",
				[]string{"11111111111x", "22222222222"},
				nil,
			},
			TestError{
				"Phones",
				errors.New(
					"the length of 12 of the value '11111111111x' is not equal to 11: not equal: validation error",
				),
			},
		},
		{
			Response{
				111,
				"",
			},
			TestError{
				"Code",
				errors.New(
					"the value '111' is not included in the list of '200,404,500': not contained: validation error",
				),
			},
		},
		{
			Response{
				200,
				"",
			},
			nil,
		},
		{
			Token{
				[]byte("xxx"),
				[]byte("zzz"),
				[]byte("yyy"),
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			errs := Validate(tt.in)

			var e *ValidationErrors
			errors.As(errs, &e)

			expected := ""
			if tt.expectedErr != nil {
				expected = tt.expectedErr.Error()
			}

			require.Equal(t, e.Error(), expected)
		})
	}
}
