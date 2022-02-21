package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (errors ValidationErrors) Error() string {
	output := ""
	for _, e := range errors {
		output += e.Field + ": " + e.Err.Error() + "\n"
	}
	return output
}

const tagName = "validate"

var ErrUnexpectedType = errors.New("unexpected type")

func Validate(v interface{}) error {
	r := reflect.ValueOf(v)
	if r.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T: %w", r, ErrUnexpectedType)
	}

	t := r.Type()
	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < t.NumField(); i++ {
		reflectStructField := t.Field(i)
		reflectValue := r.Field(i)
		fieldName := reflectStructField.Name
		validateTag := reflectStructField.Tag.Get(tagName)

		validatorTags := MakeValidatorTags(validateTag)
		if len(validatorTags) == 0 {
			continue
		}

		validationErrors = ValidateSliceField(fieldName, reflectValue, validatorTags, validationErrors)
		validationErrors = ValidateStringField(fieldName, reflectValue, validatorTags, validationErrors)
		validationErrors = ValidateIntField(fieldName, reflectValue, validatorTags, validationErrors)
	}

	fmt.Println(validationErrors)

	return validationErrors
}

func ValidateSliceField(
	fieldName string,
	reflectValue reflect.Value,
	validatorTags ValidatorTags,
	validationErrors ValidationErrors,
) ValidationErrors {
	if reflectValue.Kind() != reflect.Slice {
		return validationErrors
	}

	stringSlice, ok := reflectValue.Interface().([]string)
	if ok {
		for _, vl := range stringSlice {
			sr := reflect.ValueOf(vl)
			validationErrors = ValidateStringField(fieldName, sr, validatorTags, validationErrors)
		}
		return validationErrors
	}

	intSlice, ok := reflectValue.Interface().([]int)
	if ok {
		for _, vl := range intSlice {
			sr := reflect.ValueOf(vl)
			validationErrors = ValidateIntField(fieldName, sr, validatorTags, validationErrors)
		}
		return validationErrors
	}

	return validationErrors
}

func ValidateStringField(
	fieldName string,
	reflectValue reflect.Value,
	validatorTags ValidatorTags,
	validationErrors ValidationErrors,
) ValidationErrors {
	if reflectValue.Kind() != reflect.String {
		return validationErrors
	}

	for _, validatorTag := range validatorTags {
		validationError, err := ValidateField(fieldName, reflectValue, validatorTag)
		if err != nil {
			continue // TODO что делать с ошибками?
		}
		if validationError == nil {
			continue
		}
		validationErrors = append(validationErrors, *validationError)
	}

	return validationErrors
}

func ValidateIntField(
	fieldName string,
	reflectValue reflect.Value,
	validatorTags ValidatorTags,
	validationErrors ValidationErrors,
) ValidationErrors {
	if reflectValue.Kind() != reflect.Int {
		return validationErrors
	}

	for _, validatorTag := range validatorTags {
		validationError, err := ValidateField(fieldName, reflectValue, validatorTag)
		if err != nil {
			continue // TODO что делать с ошибками?
		}
		if validationError == nil {
			continue
		}
		validationErrors = append(validationErrors, *validationError)
	}

	return validationErrors
}

func ValidateField(fieldName string, reflectValue reflect.Value, validatorTag ValidatorTag) (*ValidationError, error) {
	name, _ := validatorTag.getName()
	parameter, _ := validatorTag.getParameter()

	validator, err := NewValidator(name)
	if err != nil {
		return nil, err
	}

	validationError := validator(reflectValue, parameter)
	if validationError == nil {
		return nil, nil
	}

	if !errors.Is(validationError, ErrValidate) {
		return nil, validationError
	}

	vErr := new(ValidationError)
	vErr.Field = fieldName
	vErr.Err = validationError
	return vErr, nil
}
