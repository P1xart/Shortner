package service

import "errors"

var (
	ErrLinkExists    = errors.New("link already exists")
	ErrLinkNotFound  = errors.New("link not found")
)
