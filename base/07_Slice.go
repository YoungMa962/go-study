package main

import "fmt"

/*
** slice 的声明和使用
 */
func main() {
	// 直接声明并赋予初始值
	slice01 := [3]int{1, 2, 3}
	fmt.Printf("slice01:[length = %d,cap = %d,value =  %v]\n", len(slice01), cap(slice01), slice01)
	// 通过 make初始化一个有初始长度的
	slice02 := make([]int, 3)
	fmt.Printf("slice02:[length = %d,,cap = %d,value =  %v]\n", len(slice02), cap(slice02), slice02)
	// slice 的判空
	var slice03 []int
	if slice03 == nil {
		fmt.Printf("slice03:[length = %d,cap = %d,value =  %v]\n", len(slice03), cap(slice03), slice03)
	}
	// slice 的遍历
	for index, value := range slice01 {
		fmt.Printf("[index = %d,value = %d]\n", index, value)
	}

	/*	append
		cap 扩容机制
		initCap = length
		newCap = oldCap * 2
	*/
	var slice04 []int
	for i := 0; i < 10; i++ {
		//fmt.Printf("slice04:[length = %d,cap = %d,value =  %v]\n", len(slice04), cap(slice04), slice04)
		slice04 = append(slice04, i)
	}
	fmt.Printf("slice04:[length = %d,cap = %d,value =  %v]\n", len(slice04), cap(slice04), slice04)

	//
	/* sub slice
		[a,b)
		[0,b)
	 	[b,len(slice))
	*/
	// [a,b)
	slice04 = slice04[0:9]
	fmt.Printf("slice04:[length = %d,cap = %d,value =  %v]\n", len(slice04), cap(slice04), slice04)
	// [0,b)
	slice04 = slice04[:5]
	fmt.Printf("slice04:[length = %d,cap = %d,value =  %v]\n", len(slice04), cap(slice04), slice04)
	// [b,len(slice))
	slice04 = slice04[3:]
	fmt.Printf("slice04:[length = %d,cap = %d,value =  %v]\n", len(slice04), cap(slice04), slice04)

	/* copy
	由于截取会改变远slice 可使用copy 进行拷贝
	copy 做的就是 将原始的slice按顺序复制到一个已经初始化完成的新slice
	容量取决于新slice的初始化
	*/
	var slice05 []int
	for i := 0; i < 10; i++ {
		slice05 = append(slice05, i)
	}
	fmt.Printf("slice05:[length = %d,cap = %d,value =  %v]\n", len(slice05), cap(slice05), slice05)

	slice05Copy := make([]int, len(slice05), cap(slice05))
	copy(slice05Copy, slice05)
	fmt.Printf("slice05Copy:[length = %d,cap = %d,value =  %v]\n", len(slice05Copy), cap(slice05Copy), slice05Copy)

	// copy 不会自动扩容
	slice05CopyNotAutoInc := make([]int, 1, 2)
	copy(slice05CopyNotAutoInc, slice05)
	fmt.Printf("slice05CopyNotAutoInc:[length = %d,cap = %d,value =  %v]\n", len(slice05CopyNotAutoInc), cap(slice05CopyNotAutoInc), slice05CopyNotAutoInc)

	// 无法直接不初始化的状态下直接copy
	var slice05CopyError []int
	copy(slice05CopyError, slice05)
	fmt.Printf("slice05CopyError:[length = %d,cap = %d,value =  %v]\n", len(slice05CopyError), cap(slice05CopyError), slice05CopyError)

}
