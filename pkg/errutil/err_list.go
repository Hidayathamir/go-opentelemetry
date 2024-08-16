package errutil

import (
	"errors"
	"fmt"
	"strings"
)

type ErrList struct {
	list []string
}

func NewErrList() *ErrList {
	return &ErrList{list: []string{}}
}

func (e *ErrList) Add(format string, a ...any) {
	if format == "" {
		return
	}
	e.list = append(e.list, fmt.Sprintf(format, a...))
}

func (e *ErrList) AddIfErr(err error) {
	if err == nil {
		return
	}
	e.Add(err.Error())
}

func (e *ErrList) IsErr() bool {
	return len(e.list) > 0
}

type Option func(*ErrOptions)

type ErrOptions struct {
	separator string
}

func WithSeparator(separator string) Option {
	return func(o *ErrOptions) {
		o.separator = separator
	}
}

func (e *ErrList) Err(options ...Option) error {
	opts := &ErrOptions{
		separator: "\n",
	}

	for _, option := range options {
		option(opts)
	}

	return errors.New(strings.Join(e.list, opts.separator))
}
