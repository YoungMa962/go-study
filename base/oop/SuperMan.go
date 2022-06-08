package oop

import "fmt"

type SuperMan struct {
	Human
}

func NewSuperMan(human Human) *SuperMan {
	return &SuperMan{Human: human}
}

func (s *SuperMan) Fly() {
	fmt.Printf("[name:%s,age:%d] is flying\n", s.Name(), s.Age())
}

func (s *SuperMan) Eat() {
	fmt.Printf("SuperMan [name:%s,age:%d] is eating\n", s.Name(), s.Age())
}

func (s *SuperMan) Sleep() {
	fmt.Printf("SuperMan [name:%s,age:%d] is sleeping\n", s.Name(), s.Age())
}
