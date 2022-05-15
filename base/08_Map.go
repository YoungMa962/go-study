package main

import (
	"fmt"
	"strconv"
)

func main() {

	var map01 map[int]string
	if map01 == nil {
		fmt.Printf("map01 is empty\n")
	}
	// map01 := make(map[int]string, 10)
	map01 = make(map[int]string)
	initMap(map01)
	print(map01)

	// modify
	map01[0] = "BEIJING_M"
	print(map01)

	//delete
	delete(map01, 4)
	print(map01)

}

func initMap(mapStruct map[int]string) {
	for i := 0; i < 10; i++ {
		mapStruct[i] = "BEIJING" + strconv.Itoa(i)
	}
}

func print(mapStruct map[int]string) {
	fmt.Println("=================================================")
	fmt.Printf("mapStruct length:%d\n", len(mapStruct))
	for key, value := range mapStruct {
		fmt.Printf("mapStruct：[key:%d , value:%s]\n", key, value)
	}
}
