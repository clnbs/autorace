package main

import (
	"fmt"
	"github.com/clnbs/autorace/internal/pkg/container"
)

func main() {
	err := container.CreateDynamicServer("mlqksjdfmlqksdfhyhakejbrf")
	if err != nil {
		fmt.Println("error while creating container :", err)
	}
}
