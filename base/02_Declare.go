package main

import "fmt"

// 全局变量声明(不可使用:= 进行全局变量声明)
var va int = 100
var vb = 100

func main() {
	// 声明 无初始值变量
	var a int
	fmt.Printf("value = %d,type = %T\n", a, a)

	// 声明 有初始值变量
	var b int = 100
	fmt.Printf("value = %d,type = %T\n", b, b)

	// 声明 变量不指定类型
	var c = 100
	fmt.Printf("value = %d,type = %T\n", c, c)

	// 声明 变量不适用var
	d := 100
	fmt.Printf("value = %d,type = %T\n", d, d)

	fmt.Printf("value = %d\t%d\n", va, vb)

	// 声明 多变量
	var aa, bb = 100, "hello"
	fmt.Printf("value = %d,type = %T\n", aa, aa)
	fmt.Printf("value = %s,type = %T\n", bb, bb)

	var (
		cc = 100
		dd = "golang"
	)
	fmt.Printf("value = %d,type = %T\n", cc, cc)
	fmt.Printf("value = %s,type = %T\n", dd, dd)

}
