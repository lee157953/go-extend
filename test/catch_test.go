package test

import (
	"errors"
	"fmt"
	. "github.com/lee157953/go-extend/ext"
	"testing"
)

var _Err = errors.New("my err")

func TestCatch(t *testing.T) {
	err := CatchErr(func() {
		Try(Foo0())
		a := Try1(Foo1())
		println(a)
		c, d := Try2(Foo2())
		println(c, d)
	})
	if err != nil {
		println(err.Error())
	}
}

func Foo0() error {
	return _Err
}

func Foo1() (string, error) {
	return "", _Err
}

func Foo2() (int, string, error) {
	return 0, "", _Err
}

func CatchErr(fn func()) (err error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
			break
		case error:
			err = r
		default:
			err = errors.New(fmt.Sprintf("unknown err: %#v", r))
		}
	}()
	fn()
	return
}
