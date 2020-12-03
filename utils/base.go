package utils

import "reflect"

/*
	参数v必须是chan, func, interface, map, pointer, or slice，否则会panic。
*/
func IsNotNil(v interface{}) bool {
	return !IsNil(v)
}

func IsNil(i interface{}) bool {
	ret := i == nil
	if !ret { //需要进一步做判断
		defer func() {
			recover()
		}()
		ret = reflect.ValueOf(i).IsNil() //值类型做异常判断，会panic的
	}
	return ret
}
