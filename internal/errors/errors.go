// Copyright 2020 Kentaro Hibino. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Package errors defines the error type and functions used by
// asynq and its internal packages.
package errors

import (
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	Code Code
	Op   Op
	Err  error
}

func (e *Error) Error() string {
	var b strings.Builder
	if e.Op != "" {
		b.WriteString(string(e.Op))
	}
	if e.Code != Unspecified {
		if b.Len() > 0 {
			b.WriteString(": ")
		}
		b.WriteString(e.Code.String())
	}
	if e.Err != nil {
		if b.Len() > 0 {
			b.WriteString(": ")
		}
		b.WriteString(e.Err.Error())
	}
	return b.String()
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Code defines the canonical error code.
type Code uint8

// List of canonical error codes.
const (
	Unspecified Code = iota
	NotFound
	FailedPrecondition
	Internal
	AlreadyExists
	Unknown
)

func (c Code) String() string {
	switch c {
	case Unspecified:
		return "ERROR_CODE_UNSPECIFIED"
	case NotFound:
		return "NOT_FOUND"
	case FailedPrecondition:
		return "FAILED_PRECONDITION"
	case Internal:
		return "INTERNAL_ERROR"
	case AlreadyExists:
		return "ALREADY_EXISTS"
	}
	panic(fmt.Sprintf("unknown error code %d", c))
}

// Op describes an operation, usually as the package and method,
// such as "rdb.Enqueue".
type Op string

func E(args ...interface{}) error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Code:
			e.Code = arg
		case error:
			e.Err = arg
		case string:
			e.Err = errors.New(arg)
		}
	}
	return e
}

// CanonicalCode returns the canonical code of the given error if one is present.
// Otherwise it returns Unspecified.
func CanonicalCode(err error) Code {
	if err == nil {
		return Unspecified
	}
	e, ok := err.(*Error)
	if !ok {
		return Unspecified
	}
	if e.Code == Unspecified {
		return CanonicalCode(e.Err)
	}
	return e.Code
}

/******************************************
    Domin Specific Error Types
*******************************************/

// TaskNotFoundError indicates that a task with the given ID does not exist
// in the given queue.
type TaskNotFoundError struct {
	Queue string // queue name
	ID    string // task id
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("cannot find task with id=%s in queue %q", e.ID, e.Queue)
}

// IsTaskNotFound reports whether any error in err's chain is of type TaskNotFoundError.
func IsTaskNotFound(err error) bool {
	var target *TaskNotFoundError
	return As(err, &target)
}

// QueueNotFoundError indicates that a queue with the given name does not exist.
type QueueNotFoundError struct {
	Queue string // queue name
}

func (e *QueueNotFoundError) Error() string {
	return fmt.Sprintf("queue %q does not exist", e.Queue)
}

// IsQueueNotFound reports whether any error in err's chain is of type QueueNotFoundError.
func IsQueueNotFound(err error) bool {
	var target *QueueNotFoundError
	return As(err, &target)
}

// TaskAlreadyArchivedError indicates that the task in question is already archived.
type TaskAlreadyArchivedError struct {
	Queue string // queue name
	ID    string // task id
}

func (e *TaskAlreadyArchivedError) Error() string {
	return fmt.Sprintf("task is already archived: id=%s, queue=%s", e.ID, e.Queue)
}

// IsTaskAlreadyArchived reports whether any error in err's chain is of type TaskAlreadyArchivedError.
func IsTaskAlreadyArchived(err error) bool {
	var target *TaskAlreadyArchivedError
	return As(err, &target)
}

/*************************************************
    Standard Library errors package functions
*************************************************/

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
//
// This function is the errors.New function from the standard libarary (https://golang.org/pkg/errors/#New).
// It is exported from this package for import convinience.
func New(text string) error { return errors.New(text) }

// Is reports whether any error in err's chain matches target.
//
// This function is the errors.Is function from the standard libarary (https://golang.org/pkg/errors/#Is).
// It is exported from this package for import convinience.
func Is(err, target error) bool { return errors.Is(err, target) }

// As finds the first error in err's chain that matches target, and if so, sets target to that error value and returns true.
// Otherwise, it returns false.
//
// This function is the errors.As function from the standard libarary (https://golang.org/pkg/errors/#As).
// It is exported from this package for import convinience.
func As(err error, target interface{}) bool { return errors.As(err, target) }

// Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// This function is the errors.Unwrap function from the standard libarary (https://golang.org/pkg/errors/#Unwrap).
// It is exported from this package for import convinience.
func Unwrap(err error) error { return errors.Unwrap(err) }
