package tests

import (
	"fmt"
	"github.com/alydnh/go-micro-ci-common/logs"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogrusScope(t *testing.T) {
	scope := &logs.LogrusScope{Entry: logrus.New().WithField("test", "test")}
	var testError error = nil
	var testResult string
	err := scope.WithField("aaa", "bbb").Call(func(ls *logs.LogrusScope) (result interface{}, err error) {
		return "test result", nil
	}).WithField("bbb", "ccc").Then(func(last interface{}, ls *logs.LogrusScope) (result interface{}, err error) {
		testResult = last.(string)
		return nil, fmt.Errorf("test error")
	}).OnError(func(err error, ls *logs.LogrusScope) error {
		testError = err
		return err
	})

	assert.NotNil(t, testError)
	assert.Equal(t, err, testError)
	assert.Equal(t, "test result", testResult)

	err = scope.WithField("aaa", "bbb").Call(func(ls *logs.LogrusScope) (result interface{}, err error) {
		return "test result1", nil
	}).WithField("bbb", "ccc").Then(func(last interface{}, ls *logs.LogrusScope) (result interface{}, err error) {
		testResult = last.(string)
		return testResult, nil
	}).OnError(func(err error, ls *logs.LogrusScope) error {
		assert.FailNow(t, "unexpected error")
		return err
	})

	assert.Nil(t, err)
	assert.Equal(t, "test result1", testResult)
}
