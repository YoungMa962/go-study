package main

import (
	"fmt"
	"go-study/oop"
)

func main() {
	human := oop.NewHuman("MaYang", 26)
	superman := oop.NewSuperMan(*human)
	//dayAction(human)
	//dayAction(superman)

	// 类型断言
	assertType(human)
	assertType(superman)
	assertType("Other")
}

func dayAction(action oop.Action) {
	action.Eat()
	action.Sleep()
}

func assertType(arg interface{}) {
	/*	human, isHuman := arg.(*oop.Human)
		superMan, isSuperMan := arg.(*oop.SuperMan)
		if isHuman {
			fmt.Println("This is a human")
			dayAction(human)
		} else if isSuperMan {
			fmt.Println("This is a human")
			dayAction(superMan)
		} else {
			fmt.Println("Type is not support")
		}*/
	fmt.Println("================================")
	switch arg.(type) {
	case *oop.Human:
		fmt.Println("This is a human")
		human := arg.(*oop.Human)
		dayAction(human)
		break
	case *oop.SuperMan:
		fmt.Println("This is a superMan")
		superMan := arg.(*oop.SuperMan)
		dayAction(superMan)
		break
	case string:
		fmt.Println("Type is string")
		break
	default:
		fmt.Println("Type is not support")
	}

}
