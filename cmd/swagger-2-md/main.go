package main

import (
	"fmt"

	"github.com/lutomas/swagger-2-md/types"
)

func main() {
	version := types.NewVersion("swagger-2-md")
	fmt.Printf("Version: %+v\n", *version)
	fmt.Println("hello world CLI")
}
