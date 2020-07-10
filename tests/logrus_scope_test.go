package tests

import (
	"fmt"
	"github.com/alydnh/go-micro-ci-common/logs"
	"github.com/alydnh/go-micro-ci-common/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogrusScopeHandle(t *testing.T) {
	scope := &logs.LogrusScope{Entry: logrus.New().WithField("test", "test")}
	var testError error = nil
	var testResult string
	err := scope.WithField("aaa", "bbb").Handle(func(ls *logs.LogrusScope) (result interface{}, err error) {
		return "test result", nil
	}).WithField("bbb", "ccc").ThenHandle(func(last interface{}, ls *logs.LogrusScope) (result interface{}, err error) {
		testResult = last.(string)
		return nil, fmt.Errorf("test error")
	}).OnError(func(err error, ls *logs.LogrusScope) error {
		testError = err
		return err
	})

	assert.NotNil(t, testError)
	assert.Equal(t, err, testError)
	assert.Equal(t, "test result", testResult)

	err = scope.WithField("aaa", "bbb").Handle(func(ls *logs.LogrusScope) (result interface{}, err error) {
		return "test result1", nil
	}).WithField("bbb", "ccc").ThenHandle(func(last interface{}, ls *logs.LogrusScope) (result interface{}, err error) {
		testResult = last.(string)
		return testResult, nil
	}).OnError(func(err error, ls *logs.LogrusScope) error {
		assert.FailNow(t, "unexpected error")
		return err
	})

	assert.Nil(t, err)
	assert.Equal(t, "test result1", testResult)
}

func TestLogrusScopeCall(t *testing.T) {
	scope := &logs.LogrusScope{Entry: logrus.New().WithField("test", "test")}
	lc := &logrusCall{T: t}
	r := scope.Call(lc.Call1, "aaa", "bbb")
	assert.Equal(t, "aaa-bbb", r.GetResult())

	r = scope.Call(lc.Call2, "bbb", "ccc")
	assert.Equal(t, "bbb-ccc", r.GetResult())

	r = scope.Call(lc.CallError, "bbb", "ccc")
	assert.Empty(t, r.GetResult())
	assert.True(t, r.HasError())
	assert.NotNil(t, r.GetError())

	r = scope.Call(lc.Call1, "1", "2").Then(lc.ThenCall)
	assert.Equal(t, r.GetResult(), "1-2")
	_ = r.Then(lc.ThenError).OnError(func(err error, ls *logs.LogrusScope) error {
		assert.Equal(t, "ThenError", err.Error())
		return err
	})

	r = scope.Call(lc.CallPanic)
	assert.Contains(t, r.GetError().Error(), "panic")
}

func TestLogrusScopeCatch(t *testing.T) {
	scope := &logs.LogrusScope{Entry: logrus.New().WithField("test", "test")}
	mustEmpty := utils.EmptyString
	err := scope.Handle(func(ls *logs.LogrusScope) (result interface{}, err error) {
		return nil, fmt.Errorf("aaa")
	}).Catch(func(err error, ls *logs.LogrusScope) error {
		assert.Equal(t, "aaa", err.Error())
		return nil
	}).ThenHandle(func(last interface{}, ls *logs.LogrusScope) (result interface{}, err error) {
		return nil, fmt.Errorf("bbb")
	}).Catch(func(err error, ls *logs.LogrusScope) error {
		return err
	}).ThenHandle(func(last interface{}, ls *logs.LogrusScope) (result interface{}, err error) {
		mustEmpty = "notEmpty"
		return
	}).OnError(func(err error, ls *logs.LogrusScope) error {
		assert.Equal(t, "bbb", err.Error())
		return err
	})
	assert.Empty(t, mustEmpty)
	assert.NotNil(t, err)
}

type logrusCall struct {
	*testing.T
	a, b, c string
}

func (lc *logrusCall) Call1(a, b string, ls *logs.LogrusScope) string {
	assert.NotNil(lc.T, ls)
	lc.a = a
	lc.b = b
	return fmt.Sprintf("%s-%s", a, b)
}

func (lc *logrusCall) Call2(a string, ls *logs.LogrusScope, b string) string {
	assert.NotNil(lc.T, ls)
	lc.a = a
	lc.b = b
	return fmt.Sprintf("%s-%s", a, b)
}

func (lc *logrusCall) ThenCall() string {
	lc.c = fmt.Sprintf("%s-%s", lc.a, lc.b)
	return lc.c
}

func (lc *logrusCall) ThenError() error {
	return fmt.Errorf("ThenError")
}

func (lc logrusCall) CallError(a string, ls *logs.LogrusScope, b string) (string, error) {
	assert.NotNil(lc.T, ls)
	return utils.EmptyString, fmt.Errorf(fmt.Sprintf("%s-%s", a, b))
}

func (lc logrusCall) CallPanic(ls *logs.LogrusScope) {
	var a *logrusCall = nil
	a.c = "aaa"
}
