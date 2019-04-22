package main

import (
	"fmt"
)

const (
	first = 1 << (10 * iota)
	second
	third
)

func main() {
	fmt.Println(first)
	fmt.Println(second)
	fmt.Println(third)
	myarray := [...]int{45, 4, 9}
	fmt.Println(myarray)
	myslice := myarray[:]
	myslice = append(myslice, 100)

	fmt.Println(myslice)
}
