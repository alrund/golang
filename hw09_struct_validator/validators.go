package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrInvalidValidator = errors.New("invalid validator for this type of field")
	ErrUnknownValidator = errors.New("unknown validator")

	ErrValid       = errors.New("validation error")
	ErrValidLength = fmt.Errorf("length exceeded: %w", ErrValid)
	ErrValidMin    = fmt.Errorf("min value exceeded: %w", ErrValid)
	ErrValidMax    = fmt.Errorf("max value exceeded: %w", ErrValid)
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
	}

	return nil, fmt.Errorf("%s: %w", name, ErrUnknownValidator)
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
			ErrValidLength)
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
		return fmt.Errorf("%d less then %d: %w", value, limit, ErrValidMin)
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
		return fmt.Errorf("%d more then %d: %w", value, limit, ErrValidMax)
	}

	return nil
}
