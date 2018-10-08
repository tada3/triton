package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tada3/triton/weather/owm"
)

const ()

func main() {
	log.Println("Triton Admin Tool")

	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	if len(args) == 0 {
		fmt.Println("Usage: att <filepath>")
		return
	}

	filepath := args[0]
	fmt.Printf("Load %v\n", filepath)
	count, err := owm.LoadCityList(filepath)
	if err != nil {
		fmt.Printf("Failed to load %v: %s\n", filepath, err.Error())
		return
	}
	fmt.Printf("count: %v\n", count)
}
