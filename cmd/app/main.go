package main

import (
	"fmt"
	"unsafe"

	"github.com/goqueue/pkg/queue"
)

func main() {
	fmt.Println("hello world.")
	fmt.Println(unsafe.Sizeof(queue.Task{}))
	fmt.Println(unsafe.Alignof(queue.Task{}))
}
