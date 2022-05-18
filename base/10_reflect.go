package main

import (
	"fmt"
	"go-study/oop"
	"reflect"
)

func main() {
	// 基本类型
	//var pai float32 = 3.1415926
	//printReflectInfo(pai)
	person := oop.NewHuman("MaYang", 26)
	//printReflectInfo(*person)
	//
	//changeValue(&pai)
	//printReflectInfo(pai)
	//
	//reflectStructField(*person)
	reflectStructMethod(*person)

}

func changeValue(arg interface{}) {
	value := reflect.ValueOf(arg)
	if value.Elem().CanSet() {
		_, ok := arg.(*float32)
		if ok {
			value.Elem().SetFloat(31.415926)
		}
	}
}

/*
 * 获取结构体 field
 */
func reflectStructField(itf interface{}) {
	rvalue := reflect.ValueOf(itf)
	fmt.Println("=======================================================")
	fmt.Printf("struct : %v\n", rvalue.Type())
	for i := 0; i < rvalue.NumField(); i++ {
		eleValue := rvalue.Field(i)
		fmt.Printf("field [value:%v ,type:%v, kind:%v]\n", eleValue, eleValue.Type(), eleValue.Kind())
	}
}

/*
 * 获取结构体 method
 */
func reflectStructMethod(itf interface{}) {
	rType := reflect.TypeOf(itf)
	fmt.Println("=======================================================")
	fmt.Printf("struct : %v\n", rType)
	for i := 0; i < rType.NumMethod(); i++ {
		method := rType.Method(i)
		fmt.Printf("[name:%s, type:%v]\n", method.Name, method.Type)
	}
}

func printReflectInfo(arg interface{}) {
	fmt.Println("=======================================================")
	var value reflect.Value
	value = reflect.ValueOf(arg)
	var typeOf reflect.Type
	typeOf = reflect.TypeOf(arg)
	var kind reflect.Kind
	kind = value.Kind()
	fmt.Printf("element [value:%v ,type:%v, kind:%v]\n", value, typeOf, kind)
}
