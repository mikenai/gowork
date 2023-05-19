package models

import (
	"errors"
	"fmt"
)

var (
	InvalidErr                    = errors.New("invalid argument")
	UserCreateParamInvalidNameErr = fmt.Errorf("invalid name: %w", InvalidErr)
	NotFoundErr                   = errors.New("not found")
)
