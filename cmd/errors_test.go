/*
 * @Author: Jeffrey <zhifeng172@163.com>
 * @Date: 2021-07-29 19:08:19
 * @Descripttion:
 */
package cmd

import (
	"serv/core/errorx"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {

	err := BadRequest.New("an_error")
	errWithContext := errorx.AddErrorContext(err, "the field is empty")
	expectedContext := "the field is empty"
	t.Log(errorx.GetErrorContext(errWithContext))

	assert.Equal(t, BadRequest, errorx.GetType(errWithContext))
	assert.Equal(t, expectedContext, errorx.GetErrorContext(errWithContext))
	assert.Equal(t, err.Error(), errWithContext.Error())
}

func TestContextInNoTypeError(t *testing.T) {
	err := errorx.New("a custom error")

	errWithContext := errorx.AddErrorContext(err, "the field is empty")
	expectedContext := "the field is empty"

	assert.Equal(t, errorx.NoType, errorx.GetType(errWithContext))
	assert.Equal(t, expectedContext, errorx.GetErrorContext(errWithContext))
	assert.Equal(t, err.Error(), errWithContext.Error())
}

func TestWrapf(t *testing.T) {
	err := errorx.New("an_error")
	wrappedError := BadRequest.Wrapf(err, "error %s", "1")

	assert.Equal(t, BadRequest, errorx.GetType(wrappedError))
	assert.EqualError(t, wrappedError, "error 1: an_error")
}

func TestWrapfInNoTypeError(t *testing.T) {
	err := errorx.Newf("an_error %s", "2")
	wrappedError := errorx.Wrapf(err, "error %s", "1")

	assert.Equal(t, errorx.NoType, errorx.GetType(wrappedError))
	assert.EqualError(t, wrappedError, "error 1: an_error 2")
}
