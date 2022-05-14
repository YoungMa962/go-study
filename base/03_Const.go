package main

import "fmt"

// 常量与枚举
func main() {
	// 常量
	const a = 100
	fmt.Printf("const value a = %d\n", a)

	//枚举 通过 iota
	const (
		BEIJING = iota
		SHANGHAI
		HANGZHOU
		NANJING
	)
	fmt.Printf("BEIJING = %d\n", BEIJING)
	fmt.Printf("SHANGHAI = %d\n", SHANGHAI)
	fmt.Printf("HANGZHOU = %d\n", HANGZHOU)
	fmt.Printf("NANJING = %d\n", NANJING)
	//iota 只可以在const中使用
	//var iotaValue int = iota
}
