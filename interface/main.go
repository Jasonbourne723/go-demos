package main

import "fmt"

func main() {

	s := SPrint(1.1, 2, 78, "nihao", nil)

	fmt.Println(s)

	// interface 内部实现 是一个包含 *类型+*数据 的元组
	// 当interface 有方法表时，还会存在一个 *itab 字段保存方法表
	var x *int = nil
	var y interface{} = x
	fmt.Println(y == nil)

	// 类型断言
	var x1 interface{} = int(10)

	// 运行时判断 “动态类型”是否是 int
	fmt.Printf("x1.(int) is %d\n", x1.(int))

	// 编译期间期间就会报错，因为go的类型转换只能发生在 “确定的静态类型”之间，而interface{}是动态的
	//fmt.Printf("int(x1) is %d", int(x1))

	x1 = float64(10)
	fmt.Println(x1.(int))
	if val, ok := x1.(float64); ok {
		fmt.Printf("x1 的类型断言正确,val=%v\n", val)
	} else {
		fmt.Println("x1 的类型断言错误")
	}

	// x1.(type) 只能在switch语句使用
	switch x1.(type) {
	case int:
		fmt.Println("x1`s type is int")
	case int64:
		fmt.Println("x1`s type is int64")
	default:
		fmt.Println("type is not support")

	}

	// 类型转换
	var a, b int32 = 10, 11
	var c int64 = 34242929424242

	// 可以转换，但是当值超出int32范围会被截断
	a = int32(c)
	fmt.Printf("%v,%v,%v\n", a, b, c)

}
