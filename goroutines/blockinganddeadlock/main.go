package main

import (
	"fmt"
)

func main() {

	channel0 := make(chan int)
	go func() {
		fmt.Println("before writting 1 to channel0...")
		channel0 <- 1
		fmt.Println("1 was written to channel0")
		fmt.Println("before writting 2 to channel0...")
		channel0 <- 2
		fmt.Println("2 was written to channel0")
	}()

	go func() {
		fmt.Printf("%d was read from channel0\n", <-channel0)
	}()

	fmt.Scanln()

	channel1 := make(chan int)

	go func() {
		fmt.Println("before writting 1 to channel1...")
		channel1 <- 1
		fmt.Println("1 was written to channel1")
		fmt.Println("before writting 2 to channel1...")
		channel1 <- 2
		fmt.Println("2 was written to channel1")
	}()

	fmt.Printf("%d was read from channel0\n", <-channel1)
	fmt.Scanln()
	fmt.Printf("%d was read from channel0\n", <-channel1)
	fmt.Scanln()

	channel2 := make(chan int)
	go func() {
		fmt.Println("before writting 1 to channel...")
		channel2 <- 1
		fmt.Println("1 was written to channel")
	}()

	fmt.Printf("%d was read from channel\n", <-channel2)
	fmt.Scanln()
	// will wait indefinitely, so should trigger a deadlock message
	fmt.Printf("%d was read from channel\n", <-channel2)
	fmt.Scanln()

}
