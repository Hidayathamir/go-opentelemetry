package errutil

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Wrap return wrapped error with msg.
func Wrap(err error, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	if err == nil {
		return errors.New(msg)
	}
	return fmt.Errorf("%s:: %w", msg, err)
}

// WrapErrMsg return wrapped err1 with err2.
// WrapErrMsg will set err2 as error message can be extract using GetMessage func.
func WrapErrMsg(err1, err2 error) error {
	if err1 == nil && err2 == nil {
		return nil
	}
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	return fmt.Errorf("~~%w~~:: %w", err2, err1)
}

// UnwrapToList unwraps error to list string.
func UnwrapToList(err error) []string {
	if err == nil {
		return nil
	}
	msg := err.Error()
	return strings.Split(msg, ":: ")
}

// GetMessage get message detail from error.
func GetMessage(err error) string {
	if err == nil {
		return ""
	}

	re := regexp.MustCompile("~~(.*?)~~")
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) > 1 {
		return matches[1]
	}

	errMsgList := strings.Split(err.Error(), ":: ")
	return errMsgList[len(errMsgList)-1]
}
