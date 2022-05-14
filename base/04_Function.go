package main

import "fmt"

/*
	函数的定义
	多返回值函数定义
*/
func main() {
	fmt.Printf("func01 = %d\n", func01(1, 1))

	sum, sub := func02(1, 1)
	fmt.Printf("func02 = [sum = %d,sub = %d]\n", sum, sub)

	sum, sub = func03(2, 1)
	fmt.Printf("func03= [sum = %d,sub = %d]\n", sum, sub)

	value, result := func04(2, 1)
	fmt.Printf("func04= [value = %d,result = %s]\n", value, result)
}

// 简单函数定义
func func01(a, b int) int {
	return a + b
}

// 多匿名相同返回值
func func02(a, b int) (int, int) {
	return a + b, a - b
}

// 多命名相同返回值
func func03(a, b int) (res1 int, res2 int) {
	fmt.Printf("res1 = %d,res2 = %d\n", res1, res2)
	res1 = a + b
	res2 = a - b
	return res1, res2
}

// 多匿名不同类型返回值
func func04(a, b int) (int, string) {
	return a + b, "success"
}
