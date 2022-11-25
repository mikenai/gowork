package models

import (
	"errors"
	"fmt"
)

var InvalidErr = errors.New("invalid argument")

var UserCreateParamInvalidNameErr = fmt.Errorf("invalid name: %w", InvalidErr)
