package errutil

import "github.com/erajayatech/go-opentelemetry/pkg/caller"

// OptAddFuncName represents options for the AddStack function.
type OptAddFuncName struct {
	Skip int
}

// OptionAddFuncName represents an option for the AddStack function.
type OptionAddFuncName func(*OptAddFuncName)

const defaultSkip = 1

// AddFuncName adds the given error with stack (the name of the calling function).
func AddFuncName(err error, options ...OptionAddFuncName) error {
	option := &OptAddFuncName{
		Skip: defaultSkip,
	}
	for _, opt := range options {
		opt(option)
	}
	return Wrap(err, caller.FuncName(caller.WithSkip(option.Skip)))
}

// WithSkip sets the number of stack frames to skip when identifying the caller.
func WithSkip(skip int) OptionAddFuncName {
	return func(wo *OptAddFuncName) {
		wo.Skip = skip + defaultSkip
	}
}
