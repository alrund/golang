package hw09structvalidator

import (
	"errors"
	"strings"
)

const validatorSeparator = "|"
const validatorNameSeparator = ":"

var (
	ErrValidatorTagNameSeparatorMissing = errors.New("validator name separator is missing")
	ErrValidatorTagNameEmpty            = errors.New("validator name is empty")
	ErrValidatorTagParameterEmpty       = errors.New("validator parameter is empty")
)

type ValidatorTags []ValidatorTag

type ValidatorTag string

func (v ValidatorTag) getName() (string, error) {
	str := string(v)
	index := strings.IndexAny(str, validatorNameSeparator)
	if index < 0 {
		return "", ErrValidatorTagNameSeparatorMissing
	}
	name := str[:index]
	if len(name) == 0 {
		return "", ErrValidatorTagNameEmpty
	}
	return name, nil
}

func (v ValidatorTag) getParameter() (string, error) {
	name, err := v.getName()
	if err != nil {
		return "", err
	}
	parameter := string(v[len(name)+1:])
	if len(parameter) == 0 {
		return "", ErrValidatorTagParameterEmpty
	}
	return parameter, nil
}

func MakeValidatorTags(tag string) ValidatorTags {
	if len(tag) == 0 {
		return make(ValidatorTags, 0)
	}

	values := strings.Split(tag, validatorSeparator)
	vt := make(ValidatorTags, 0, len(values))

	for _, value := range values {
		vt = append(vt, ValidatorTag(value))
	}

	return vt
}
