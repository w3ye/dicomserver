package service

import (
	"errors"
	"strconv"

	"github.com/suyashkumar/dicom/pkg/tag"
)

type queryTag string

func (q queryTag) validateFormat() error {
	// tag query should be in the format of (XXXX,XXXX)
	// total length should be 11
	if len(q) != 11 {
		return errors.New("invalid tag length")
	}

	// check if the first element is wrapped by paranthesis
	if q[0] != '(' || q[len(q)-1] != ')' {
		return errors.New("invalid tag format")
	}

	return nil
}

func (q queryTag) extractGroupAndElement() (uint16, uint16, error) {
	if err := q.validateFormat(); err != nil {
		return 0, 0, err
	}
	group, err := strconv.ParseUint(string(q[1:5]), 16, 0)
	if err != nil {
		return 0, 0, err
	}

	element, err := strconv.ParseUint(string(q[6:10]), 16, 0)
	if err != nil {
		return 0, 0, err
	}

	return uint16(group), uint16(element), nil
}

func (q queryTag) convertQueryTagToDicomTag() (tag.Tag, error) {
	groupNumber, elementNumber, err := q.extractGroupAndElement()
	if err != nil {
		return tag.Tag{}, err
	}

	tag := tag.Tag{
		Group:   groupNumber,
		Element: elementNumber,
	}

	return tag, nil
}
