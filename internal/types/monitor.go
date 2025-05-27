package types

import (
	"fmt"
)

type ErrorType int

const (
	ClientError ErrorType = iota
	ConfigError
)

type MonitorError struct {
	Type ErrorType
	Err  error
}

func (e *MonitorError) Error() string {
	if e.Type == ConfigError {
		return fmt.Sprintf("configuration error: %v", e.Err)
	}
	return fmt.Sprintf("%v", e.Err)
}

func IsClientError(err error) bool {
	if me, ok := err.(*MonitorError); ok {
		return me.Type == ClientError
	}
	return false
}

func IsConfigError(err error) bool {
	if me, ok := err.(*MonitorError); ok {
		return me.Type == ConfigError
	}
	return false
}

func NewClientError(err error) *MonitorError {
	return &MonitorError{
		Type: ClientError,
		Err:  err,
	}
}

func NewConfigError(err error) *MonitorError {
	return &MonitorError{
		Type: ConfigError,
		Err:  err,
	}
}
