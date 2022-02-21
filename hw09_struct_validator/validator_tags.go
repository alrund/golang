package hw09structvalidator

import "strings"

const validatorSeparator = "|"

type ValidatorTags struct {
	tags []ValidatorTag
}

func NewValidatorTags(tag string) *ValidatorTags {
	if tag == "" {
		return nil
	}

	vt := new(ValidatorTags)
	values := strings.Split(tag, validatorSeparator)
	for _, value := range values {
		vt.tags = append(vt.tags, ValidatorTag(value))
	}
	return vt
}
