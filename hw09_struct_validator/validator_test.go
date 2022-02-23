package hw09structvalidator

import (
	"encoding/json"
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
			ValidationErrors{{
				Field: "Age",
				Err: fmt.Errorf(
					"the length of 3 of the value '111' is not equal to 36: %w: validation error",
					ErrValidateLength,
				),
			}},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 3, "test@test.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ValidationErrors{{
				Field: "Age",
				Err: fmt.Errorf(
					"3 less then 18: %w: validation error",
					ErrValidateMin,
				),
			}},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 333, "test@test.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ValidationErrors{{
				Field: "Age",
				Err: fmt.Errorf(
					"333 more then 50: %w: validation error",
					ErrValidateMax,
				),
			}},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "testtest.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ValidationErrors{{
				Field: "Email",
				Err: fmt.Errorf(
					"the 'testtest.ru' value does not match the '^\\w+@\\w+\\.\\w+$' pattern: %w: validation error",
					ErrValidateRegexp,
				),
			}},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "test@test.ru", "not",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ValidationErrors{{
				Field: "Role",
				Err: fmt.Errorf(
					"the value 'not' is not included in the list of 'admin,stuff': %w: validation error",
					ErrValidateIn,
				),
			}},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "test@test.ru", "stuff",
				[]string{"11111111111x", "22222222222"},
				nil,
			},
			ValidationErrors{{
				Field: "Phones",
				Err: fmt.Errorf(
					"the length of 12 of the value '11111111111x' is not equal to 11: %w: validation error",
					ErrValidateLength,
				),
			}},
		},
		{
			Response{
				111,
				"",
			},
			ValidationErrors{{
				Field: "Code",
				Err: fmt.Errorf(
					"the value '111' is not included in the list of '200,404,500': %w: validation error",
					ErrValidateIn,
				),
			}},
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
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.Nil(t, err)
			} else {
				var validationErrs ValidationErrors
				require.ErrorAs(t, err, &validationErrs)
				require.Equal(t, err.Error(), validationErrs.Error())
			}
		})
	}
}
