package models

import (
	"errors"
	"fmt"
)

var InvalidArgumentErr = errors.New("invalid argument")

var UserCreateParamInvalidNameErr = fmt.Errorf("invalid name: %w", InvalidArgumentErr)

var UserUpdateParamInvalidNameEmptyErr = fmt.Errorf("invalid name: %w", InvalidArgumentErr)

var UserUpdateParamInvalidNameLenErr = fmt.Errorf("invalid name: %w", InvalidArgumentErr)

var NotFoundErr = errors.New("not found")
