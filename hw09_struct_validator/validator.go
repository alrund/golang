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
	structReflectValue := reflect.ValueOf(v)
	if structReflectValue.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T: %w", structReflectValue, ErrUnexpectedType)
	}

	var validationErrors ValidationErrors

	structReflectType := structReflectValue.Type()
	for i := 0; i < structReflectType.NumField(); i++ {
		reflectStructField := structReflectType.Field(i)
		reflectValue := structReflectValue.Field(i)
		structTag := reflectStructField.Tag.Get(tagName)

		validatorTags := getValidatorTags(structTag)
		if len(validatorTags) == 0 {
			continue
		}

		var (
			fieldValidationErrors ValidationErrors
			err                   error
		)

		switch reflectValue.Kind() {
		case reflect.Slice:
			fieldValidationErrors, err = validateSliceField(reflectStructField.Name, reflectValue, validatorTags)
		case reflect.String, reflect.Int:
			fieldValidationErrors, err = validateSimpleField(reflectStructField.Name, reflectValue, validatorTags)
		case // for golangci-lint
			reflect.Array,
			reflect.Bool,
			reflect.Chan,
			reflect.Complex128,
			reflect.Complex64,
			reflect.Float32,
			reflect.Float64,
			reflect.Func,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Int8,
			reflect.Interface,
			reflect.Invalid,
			reflect.Map,
			reflect.Ptr,
			reflect.Struct,
			reflect.Uint,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uint8,
			reflect.Uintptr,
			reflect.UnsafePointer:
			continue
		}

		if err != nil {
			return err
		}
		validationErrors = append(validationErrors, fieldValidationErrors...)
	}

	fmt.Println(validationErrors)

	return validationErrors
}

func validateSliceField(
	fieldName string, sliceReflectValue reflect.Value, validatorTags ValidatorTags,
) (ValidationErrors, error) {
	reflectValues, ok := getSliceReflectValues(sliceReflectValue)
	if !ok {
		return nil, nil
	}

	var validationErrors ValidationErrors
	for _, reflectValue := range reflectValues {
		fieldValidationErrors, err := validateSimpleField(fieldName, reflectValue, validatorTags)
		if err != nil {
			return nil, err
		}
		validationErrors = append(validationErrors, fieldValidationErrors...)
	}
	return validationErrors, nil
}

func validateSimpleField(
	fieldName string, reflectValue reflect.Value, validatorTags ValidatorTags,
) (ValidationErrors, error) {
	validationErrors := make(ValidationErrors, 0)

	for _, validatorTag := range validatorTags {
		validationError, err := useValidator(fieldName, reflectValue, validatorTag)
		if err != nil {
			return nil, err
		}
		if validationError == nil {
			continue
		}
		validationErrors = append(validationErrors, *validationError)
	}

	return validationErrors, nil
}

func useValidator(fieldName string, reflectValue reflect.Value, validatorTag ValidatorTag) (*ValidationError, error) {
	name, _ := validatorTag.getName()
	parameter, _ := validatorTag.getParameter()

	validator, err := NewValidator(name)
	if err != nil {
		return nil, err
	}

	verr := validator(reflectValue, parameter)
	if !errors.Is(verr, ErrValidate) {
		return nil, verr
	}

	return &ValidationError{fieldName, verr}, nil
}

func getSliceReflectValues(reflectValue reflect.Value) ([]reflect.Value, bool) {
	if stringReflectValues, ok := getStringSliceReflectValues(reflectValue); ok {
		return stringReflectValues, true
	}
	if intReflectValues, ok := getIntSliceReflectValues(reflectValue); ok {
		return intReflectValues, true
	}

	return nil, false
}

func getStringSliceReflectValues(reflectValue reflect.Value) ([]reflect.Value, bool) {
	values := make([]reflect.Value, 0)
	items, ok := reflectValue.Interface().([]string)
	if !ok {
		return values, false
	}

	for _, item := range items {
		values = append(values, reflect.ValueOf(item))
	}

	return values, true
}

func getIntSliceReflectValues(reflectValue reflect.Value) ([]reflect.Value, bool) {
	values := make([]reflect.Value, 0)
	items, ok := reflectValue.Interface().([]int)
	if !ok {
		return values, false
	}

	for _, item := range items {
		values = append(values, reflect.ValueOf(item))
	}

	return values, true
}
