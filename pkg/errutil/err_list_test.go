package errutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrList(t *testing.T) {
	t.Run("add 0 time", func(t *testing.T) {
		e := NewErrList()
		assert.False(t, e.IsErr())
		assert.Equal(t, "", e.Err().Error())
	})
	t.Run("add 1 time empty string", func(t *testing.T) {
		e := NewErrList()
		e.Add("")
		assert.False(t, e.IsErr())
		assert.Equal(t, "", e.Err().Error())
	})
	t.Run("add 1 time", func(t *testing.T) {
		e := NewErrList()
		e.Add("err1")
		assert.True(t, e.IsErr())
		assert.Equal(t, "err1", e.Err().Error())
	})
	t.Run("add 2 time", func(t *testing.T) {
		e := NewErrList()
		e.Add("err1")
		e.Add("err2")
		assert.True(t, e.IsErr())
		assert.Equal(t, "err1\nerr2", e.Err().Error())
	})
	t.Run("add 2 time with separator", func(t *testing.T) {
		e := NewErrList()
		e.Add("err1")
		e.Add("err2")
		assert.True(t, e.IsErr())
		assert.Equal(t, "err1\n\nerr2", e.Err(WithSeparator("\n\n")).Error())
	})
	t.Run("add err", func(t *testing.T) {
		e := NewErrList()
		e.AddIfErr(errors.New("hehe"))
		assert.True(t, e.IsErr())
		assert.Equal(t, "hehe", e.Err().Error())
	})
	t.Run("add nil err", func(t *testing.T) {
		e := NewErrList()
		e.AddIfErr(nil)
		assert.False(t, e.IsErr())
		assert.Equal(t, "", e.Err().Error())
	})
}
