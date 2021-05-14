package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

var role *UserRole

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in:          nil,
		expectedErr: ErrValueIsNotStruct,
	},
	{
		in:          role,
		expectedErr: ErrValueIsNotStruct,
	},
	{
		in:          true,
		expectedErr: ErrValueIsNotStruct,
	},
	{
		in: struct {
			Guest bool `validate:"equal:true"`
		}{},
		expectedErr: errors.New("unable to validate 'Guest' field: field type is not supported"),
	},
	{
		in: struct {
			Value int `validate:"bigint:true"`
		}{},
		expectedErr: errors.New("unable to validate 'Value' field with validator 'bigint': validator type is not supported"),
	},
	{
		in: struct {
			Value int `validate:"min"`
		}{},
		expectedErr: errors.New("unable to validate 'Value' field with validator 'min': validator has wrong settings"),
	},
	{
		in: struct {
			Age int `validate:"min:x"`
		}{},
		expectedErr: errors.New("unable to validate 'Age' field with validator 'min': strconv.Atoi: parsing \"x\": invalid syntax"),
	},
	{
		in: struct {
			Age int `validate:"max:y"`
		}{},
		expectedErr: errors.New("unable to validate 'Age' field with validator 'max': strconv.Atoi: parsing \"y\": invalid syntax"),
	},
	{
		in: struct {
			Age int `validate:"in:z,y,x"`
		}{},
		expectedErr: errors.New("unable to validate 'Age' field with validator 'in': strconv.Atoi: parsing \"z\": invalid syntax"),
	},
	{
		in: struct {
			Name string `validate:"email:true"`
		}{},
		expectedErr: errors.New("unable to validate 'Name' field with validator 'email': validator type is not supported"),
	},
	{
		in: struct {
			Name string `validate:"len"`
		}{},
		expectedErr: errors.New("unable to validate 'Name' field with validator 'len': validator has wrong settings"),
	},
	{
		in: struct {
			Name string `validate:"len:x"`
		}{},
		expectedErr: errors.New("unable to validate 'Name' field with validator 'len': strconv.Atoi: parsing \"x\": invalid syntax"),
	},
	{
		in: struct {
			Name string `validate:"regexp:\\h+"`
		}{},
		expectedErr: errors.New("unable to validate 'Name' field with validator 'regexp': error parsing regexp: invalid escape sequence: `\\h`"),
	},
	{
		in: Token{},
	},
	{
		in:          App{},
		expectedErr: ValidationErrors([]ValidationError{{"Version", errors.New("must be 5 characters long")}}),
	},
	{
		in: App{"0.0.1"},
	},
	{
		in:          Response{Code: 204},
		expectedErr: ValidationErrors([]ValidationError{{"Code", errors.New("must be one of [200 404 500]")}}),
	},
	{
		in: Response{Code: 404},
	},
	{
		in: User{
			Phones: []string{"x", "y"},
		},
		expectedErr: ValidationErrors([]ValidationError{
			{"ID", errors.New("must be 36 characters long")},
			{"Age", errors.New("must be greater than 18")},
			{"Email", errors.New("must satisfy regular expression '^\\w+@\\w+\\.\\w+$'")},
			{"Role", errors.New("must be one of [\"admin\" \"stuff\"]")},
			{"Phones_0", errors.New("must be 11 characters long")},
			{"Phones_1", errors.New("must be 11 characters long")},
		}),
	},
	{
		in: User{
			Age:    51,
			meta:   []byte{},
			Phones: []string{strings.Repeat("x", 11), "y"},
		},
		expectedErr: ValidationErrors([]ValidationError{
			{"ID", errors.New("must be 36 characters long")},
			{"Age", errors.New("must be less than 50")},
			{"Email", errors.New("must satisfy regular expression '^\\w+@\\w+\\.\\w+$'")},
			{"Role", errors.New("must be one of [\"admin\" \"stuff\"]")},
			{"Phones_1", errors.New("must be 11 characters long")},
		}),
	},
	{
		in: User{
			Age:    35,
			Role:   "admin",
			Email:  "john@otus.ru",
			ID:     strings.Repeat("x", 36),
			Phones: []string{strings.Repeat("x", 11), strings.Repeat("y", 11)},
		},
	},
	{
		in: struct {
			Res Response `validate:"nested"`
		}{
			Res: Response{Code: 204},
		},
		expectedErr: ValidationErrors([]ValidationError{
			{"Res", ValidationErrors([]ValidationError{{"Code", errors.New("must be one of [200 404 500]")}})},
		}),
	},
	{
		in: struct {
			Res Response `validate:"nested"`
		}{
			Res: Response{Code: 404},
		},
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			if tt.expectedErr != nil {
				require.EqualError(t, Validate(tt.in), tt.expectedErr.Error())
			} else {
				require.NoError(t, Validate(tt.in))
			}
		})
	}
}
