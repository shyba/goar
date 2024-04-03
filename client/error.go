package client

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrPendingTx    = errors.New("pending")
	ErrInvalidId    = errors.New("invalid arweave id")
	ErrBadGateway   = errors.New("bad gateway")
	ErrRequestLimit = errors.New("arweave gateway request limit")
)
