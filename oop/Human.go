package oop

import (
	"fmt"
)

type Human struct {
	name string
	age  int
}

func NewHuman(name string, age int) *Human {
	return &Human{name: name, age: age}
}

func (h *Human) Age() int {
	return h.age
}

func (h *Human) SetAge(age int) {
	h.age = age
}

func (h *Human) Name() string {
	return h.name
}

func (h *Human) SetName(name string) {
	h.name = name
}

func (h *Human) Eat() {
	fmt.Printf("Human [name:%s,age:%d] is eating\n", h.Name(), h.Age())
}

func (h *Human) Sleep() {
	fmt.Printf("Human [name:%s,age:%d] is sleeping\n", h.Name(), h.Age())
}
