package main

import (
	"fmt"
	"time"
)

func main() {
	// timeout for each channel read
	c1 := boring("step1")
	func() {
		for {
			select {
			case s := <-c1:
				fmt.Println(s)
			case <-time.After(time.Second):
				fmt.Println("Pt1 Timed out")
				return
			}
		}
	}()

	// timeout for overall
	c2 := boring("step2")
	func() {
		timeout := time.After(2 * time.Second)
		for {
			select {
			case s := <-c2:
				fmt.Println(s)
			case <-timeout:
				fmt.Println("Pt2 Timed out")
				return
			}
		}
	}()
}

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 1; i <= 10; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Second)
		}
	}()
	return c
}
