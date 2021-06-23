package main

import "fmt"

type subject interface{
	 add()()
}

type realSubject struct{}

//implements subject interface
func (r realSubject)add()(){
	fmt.Println("realSubject add()")
}

type proxySubject struct{
	realSubject subject
}

func (r proxySubject)add()(){
	fmt.Println("before proxy")
	r.realSubject.add()
	fmt.Print("after proxy")
}

func  NewProxySubject(subject2 realSubject)(*proxySubject){
	r := proxySubject{}
	r.realSubject = subject2
	return &r;
}
func main()(){
	var proxy subject = NewProxySubject(realSubject{})
	proxy.add()
}
