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
			ErrValidateLength,
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 3, "test@test.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ErrValidateMin,
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 333, "test@test.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ErrValidateMax,
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "testtest.ru", "stuff",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ErrValidateRegexp,
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "test@test.ru", "not",
				[]string{"11111111111", "22222222222"},
				nil,
			},
			ErrValidateIn,
		},
		{
			User{
				"123456789012345678901234567890123456",
				"xxx", 33, "test@test.ru", "stuff",
				[]string{"11111111111x", "22222222222"},
				nil,
			},
			ErrValidateLength,
		},
		{
			Response{
				111,
				"",
			},
			ErrValidateIn,
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
			for _, v := range errs.(ValidationErrors) {
				require.ErrorIs(t, v.Err, tt.expectedErr)
			}
		})
	}
}
