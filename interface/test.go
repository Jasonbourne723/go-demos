package main

import (
	"reflect"
	"strconv"
	"strings"
)

func SPrint(args ...interface{}) string {

	if len(args) == 0 {
		return ""
	}
	ans := strings.Builder{}
	for i := 0; i < len(args); i++ {

		switch v := args[i].(type) {
		case int:
			ans.WriteString(strconv.Itoa(v))
		case int32:
			ans.WriteString(strconv.Itoa(int(v)))
		case int64:
			ans.WriteString(strconv.FormatInt(v, 10))
		case string:
			ans.WriteString(v)
		case nil:
			ans.WriteString("nil")
		default:
			ans.WriteString("<" + reflect.TypeOf(v).String() + "is not impletement" + ">")
		}
		if i != len(args)-1 {
			ans.WriteString(",")
		}
	}
	return ans.String()
}
