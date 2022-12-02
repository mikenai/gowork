package models

import (
	"errors"
	"fmt"
)

var InvalidErr = errors.New("invalid argument")
var NotFound = errors.New("not found")
var UserCreateParamInvalidNameErr = fmt.Errorf("invalid name: %w", InvalidErr)
