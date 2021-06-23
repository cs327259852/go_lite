package main

import (
	"fmt"
	"testing"
)

//通过匿名结构体实现继承
func TestStringer(t *testing.T){
	var c ChoiceQuestion = ChoiceQuestion{}
	c.Get()
	c.name = "hello"
	c.BaseQuestion.name = "world"
	fmt.Println(c)

}

type cat struct{
	name string
	age int32
}
func (c cat)String()string{
	return string(c.age)+string(c.age)
}

type BaseQuestion struct{
	name string
	age int
}

func (self BaseQuestion)Get()(){
	fmt.Print(123)
}

func (self BaseQuestion)Post()(){}

type ChoiceQuestion struct{
	 BaseQuestion
	 name string
}

