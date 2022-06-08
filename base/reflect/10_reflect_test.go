package reflect

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type User struct {
	Name    string
	Age     int
	Married bool
}

func NewUser(name string, age int, married bool) *User {
	return &User{Name: name, Age: age, Married: married}
}

func Add(a, b int) int {
	return a + b
}

func (receiver *User) Say(words string) string {
	return receiver.Name + words
}

func TestInspectStruct(t *testing.T) {
	user := NewUser("YangMa", 16, false)
	v := reflect.ValueOf(*user)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		// 布尔
		case reflect.Bool:
			fmt.Printf("field:%d type:%s value:%t\n", i, field.Type().Name(), field.Bool())
			break
		// string
		case reflect.String:
			fmt.Printf("field:%d type:%s value:%q\n", i, field.Type().Name(), field.String())
			break
		// 有符号整型
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Printf("field:%d type:%s value:%d\n", i, field.Type().Name(), field.Int())
			break
		// 无符号整型
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fmt.Printf("field:%d type:%s value:%d\n", i, field.Type().Name(), field.Uint())
			break
		// 浮点型
		case reflect.Float64, reflect.Float32:
			fmt.Printf("field:%d type:%s value:%f\n", i, field.Type().Name(), field.Float())
			break
		default:
			fmt.Printf("field:%d unhandled kind:%s\n", i, field.Kind())
		}
	}
}

func TestInspectFunction(t *testing.T) {
	v := reflect.TypeOf(Add)
	fmt.Println("函数入参")
	for i := 0; i < v.NumIn(); i++ {
		in := v.In(i)
		fmt.Printf("[Index : %d,type:%s]\n", i, in.Name())
	}
	fmt.Println("函数返回值")
	for i := 0; i < v.NumOut(); i++ {
		in := v.In(i)
		fmt.Printf("[Index : %d,type:%s]\n", i, in.Name())
	}
}

func TestInspectMethod(t *testing.T) {
	user := NewUser("YangMa", 16, false)
	typeOf := reflect.TypeOf(user)
	for i := 0; i < typeOf.NumMethod(); i++ {
		method := typeOf.Method(i)
		argv := make([]string, 0, method.Type.NumIn())
		results := make([]string, 0, method.Type.NumOut())
		for j := 1; j < method.Type.NumIn(); j++ {
			argv = append(argv, method.Type.In(j).Name())
		}
		for j := 0; j < method.Type.NumOut(); j++ {
			results = append(results, method.Type.Out(j).Name())
		}
		fmt.Printf("func (receiver *%s) %s(%s) %s\n",
			typeOf.Elem().Name(),
			method.Name,
			strings.Join(argv, ","),
			strings.Join(results, ","))
	}

}

func TestMethodCall(t *testing.T) {
	user := NewUser("YangMa", 16, false)
	invoke(user.Say, "hello")
	invokeByMethodName(user, "Say", "\tHello")
}

func invoke(method interface{}, args ...interface{}) {
	v := reflect.ValueOf(method)
	argsV := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		argsV = append(argsV, reflect.ValueOf(arg))
	}
	result := v.Call(argsV)
	fmt.Println("result:")
	for _, ret := range result {
		fmt.Println(ret.Interface())
	}
}
func invokeByMethodName(instance interface{}, methodName string, args ...interface{}) {
	v := reflect.ValueOf(instance)
	m := v.MethodByName(methodName)
	argsV := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		argsV = append(argsV, reflect.ValueOf(arg))
	}
	result := m.Call(argsV)
	fmt.Println("result:")
	for _, ret := range result {
		fmt.Println(ret.Interface())
	}
}
