package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidValidator = errors.New("invalid validator for this type of field")
	ErrUnknownValidator = errors.New("unknown validator")

	ErrValidate       = errors.New("validation error")
	ErrValidateLength = fmt.Errorf("length exceeded: %w", ErrValidate)
	ErrValidateMin    = fmt.Errorf("min value exceeded: %w", ErrValidate)
	ErrValidateMax    = fmt.Errorf("max value exceeded: %w", ErrValidate)
	ErrValidateIn     = fmt.Errorf("not contained: %w", ErrValidate)
	ErrValidateRegexp = fmt.Errorf("not match: %w", ErrValidate)
)

type Validator func(reflectValue reflect.Value, parameter string) error

func NewValidator(name string) (Validator, error) {
	switch name {
	case "len":
		return LenValidator, nil
	case "min":
		return MinValidator, nil
	case "max":
		return MaxValidator, nil
	case "regexp":
		return RegexpValidator, nil
	case "in":
		return InValidator, nil
	}

	return nil, fmt.Errorf("%s: %w", name, ErrUnknownValidator)
}

func RegexpValidator(reflectValue reflect.Value, parameter string) error {
	if reflectValue.Kind() != reflect.String {
		return ErrInvalidValidator
	}

	value := reflectValue.String()

	regExp, err := regexp.Compile(parameter)
	if err != nil {
		return err
	}

	if !regExp.MatchString(value) {
		return fmt.Errorf("the '%s' value does not match the '%s' pattern: %w", value, parameter, ErrValidateRegexp)
	}

	return nil
}

func InValidator(reflectValue reflect.Value, parameter string) error {
	var inspectedValue string
	switch reflectValue.Kind() {
	case reflect.Int:
		inspectedValue = strconv.Itoa(int(reflectValue.Int()))
	case reflect.String:
		inspectedValue = reflectValue.String()
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
		reflect.Slice,
		reflect.Struct,
		reflect.Uint,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uint8,
		reflect.Uintptr,
		reflect.UnsafePointer:
		return ErrInvalidValidator
	}

	validateResult := false
	for _, value := range strings.Split(parameter, ",") {
		if inspectedValue == strings.TrimSpace(value) {
			validateResult = true
			break
		}
	}

	if !validateResult {
		return fmt.Errorf(
			"the value '%s' is not included in the list of '%s': %w",
			inspectedValue,
			parameter,
			ErrValidateIn)
	}

	return nil
}

func LenValidator(reflectValue reflect.Value, parameter string) error {
	if reflectValue.Kind() != reflect.String {
		return ErrInvalidValidator
	}

	value := reflectValue.String()

	limit, err := strconv.Atoi(parameter)
	if err != nil {
		return err
	}

	if length := len(value); length > limit {
		return fmt.Errorf(
			"the length %d of the phone %s is greater than the limit %d: %w",
			length,
			value,
			limit,
			ErrValidateLength)
	}

	return nil
}

func MinValidator(reflectValue reflect.Value, parameter string) error {
	if reflectValue.Kind() != reflect.Int {
		return ErrInvalidValidator
	}

	limit, err := strconv.Atoi(parameter)
	if err != nil {
		return err
	}

	if value := reflectValue.Int(); value < int64(limit) {
		return fmt.Errorf("%d less then %d: %w", value, limit, ErrValidateMin)
	}

	return nil
}

func MaxValidator(reflectValue reflect.Value, parameter string) error {
	if reflectValue.Kind() != reflect.Int {
		return ErrInvalidValidator
	}

	limit, err := strconv.Atoi(parameter)
	if err != nil {
		return err
	}

	if value := reflectValue.Int(); value > int64(limit) {
		return fmt.Errorf("%d more then %d: %w", value, limit, ErrValidateMax)
	}

	return nil
}
