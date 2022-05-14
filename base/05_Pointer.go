package main

import "fmt"

// 指针
func main() {
	a, b := 10, 20
	fmt.Printf("[a = %d,b = %d]\n", a, b)
	swap(&a, &b)
	fmt.Printf("After swap [a = %d,b = %d]\n", a, b)

	// 二级指针
	var pA = &a
	fmt.Printf("[*pA = %d,a = %d]\n", *pA, a)
	var ppA = &pA
	fmt.Printf("[**ppA= %d,a = %d]\n", **ppA, a)
}

func swap(valA, valB *int) {
	fmt.Printf("swap [valA = %d,valB = %d]\n", valA, valB)
	fmt.Printf("swap [*valA= %d,*valB = %d]\n", *valA, *valB)

	var temp int
	temp = *valA
	*valA = *valB
	*valB = temp
}
