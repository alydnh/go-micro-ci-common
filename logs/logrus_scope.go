package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type LogrusScopeCallHandler func(ls *LogrusScope) (result interface{}, err error)
type LogrusScopeThenHandler func(last interface{}, ls *LogrusScope) (result interface{}, err error)
type LogrusScopeErrorHandler func(err error, ls *LogrusScope) error

type LogrusScope struct {
	*logrus.Entry
}

func (ls LogrusScope) WithFields(fields logrus.Fields) *LogrusScope {
	return &LogrusScope{ls.Entry.WithFields(fields)}
}

func (ls LogrusScope) WithField(key string, value interface{}) *LogrusScope {
	return &LogrusScope{ls.Entry.WithField(key, value)}
}

func (ls *LogrusScope) Call(h interface{}, args ...interface{}) (result *LogrusScopeResult) {
	return internalCall(ls, h, args)
}

func (ls *LogrusScope) Handle(h LogrusScopeCallHandler) *LogrusScopeResult {
	result, err := h(ls)
	return &LogrusScopeResult{
		err:    err,
		Entry:  ls.Entry,
		result: result,
	}
}

type LogrusScopeResult struct {
	err    error
	result interface{}
	*logrus.Entry
}

func (r *LogrusScopeResult) ThenHandle(h LogrusScopeThenHandler) *LogrusScopeResult {
	if r.HasError() {
		return r
	}
	result, err := h(r.result, &LogrusScope{r.Entry})
	return &LogrusScopeResult{
		err:    err,
		Entry:  r.Entry,
		result: result,
	}
}

func (r *LogrusScopeResult) Then(h interface{}, args ...interface{}) (result *LogrusScopeResult) {
	if r.HasError() {
		return r
	}
	return internalCall(&LogrusScope{r.Entry}, h, args)
}

func (r LogrusScopeResult) WithFields(fields logrus.Fields) *LogrusScopeResult {
	return &LogrusScopeResult{
		err:    r.err,
		result: r.result,
		Entry:  r.Entry.WithFields(fields),
	}
}

func (r LogrusScopeResult) WithField(key string, value interface{}) *LogrusScopeResult {
	return &LogrusScopeResult{
		err:    r.err,
		result: r.result,
		Entry:  r.Entry.WithField(key, value),
	}
}

func (r *LogrusScopeResult) HasError() bool {
	return nil != r.err
}

func (r *LogrusScopeResult) GetError() error {
	return r.err
}

func (r *LogrusScopeResult) GetResult() interface{} {
	return r.result
}

func (r *LogrusScopeResult) OnError(h LogrusScopeErrorHandler) error {
	if r.HasError() {
		return h(r.err, &LogrusScope{r.Entry})
	}

	return nil
}

func internalCall(scope *LogrusScope, h interface{}, args []interface{}) (result *LogrusScopeResult) {
	result = &LogrusScopeResult{
		Entry: scope.Entry,
	}

	doRecover := true
	defer func() {
		if doRecover {
			v := recover()
			if nil != v {
				if err, ok := v.(error); ok {
					result.err = err
				} else {
					result.err = fmt.Errorf("panic: %v", v)
				}
			}
		}
	}()

	t := reflect.TypeOf(h)
	if t.Kind() != reflect.Func {
		result.err = fmt.Errorf("h is not a function")
		return result
	}
	f := reflect.ValueOf(h)
	in := make([]reflect.Value, 0, t.NumIn())

	scopeType := reflect.TypeOf(scope)
	for i := 0; i < t.NumIn(); i++ {
		arg := t.In(i)
		if arg == scopeType {
			in = append(in, reflect.ValueOf(scope))
		} else {
			if len(args) == 0 {
				doRecover = false
				panic(fmt.Errorf("args length too small"))

			}
			in = append(in, reflect.ValueOf(args[0]))
			if len(args) > 0 {
				args = args[1:]
			}
		}
	}

	if len(args) > 0 {
		doRecover = false
		panic(fmt.Errorf("args length too large"))
	}

	outs := f.Call(in)
	outInterfaces := make([]interface{}, 0, len(outs))
	for _, out := range outs {
		v := out.Interface()
		if nil != v {
			if err, ok := v.(error); ok {
				result.err = err
				return result
			}
		}
		outInterfaces = append(outInterfaces, v)
	}
	if len(outInterfaces) == 1 {
		result.result = outInterfaces[0]
	} else if len(outInterfaces) > 0 {
		result.result = outInterfaces
	}
	return result
}
