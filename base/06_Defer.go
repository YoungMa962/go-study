package main

import "fmt"

/*
** defer 关键字 类似Java的 final
** 1、函数执行完，最后触发
** 2、和return 的执行顺序
 */
func main() {
	testProgress()
}

func testProgress() int {
	defer deferCallA()
	defer deferCallB()
	fmt.Println("function is running")
	return returnCall()
}

func returnCall() int {
	fmt.Println("returnCall function is running")
	return 0
}

func deferCallA() {
	fmt.Println("deferCallA is running")
}

func deferCallB() {
	fmt.Println("deferCallB is running")
}
