package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tada3/triton/tritondb"
	"github.com/tada3/triton/weather/owm"
)

const ()

var (
	clearFlag bool
)

func main() {
	log.Println("Triton Admin Tool")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:\n  att [--clear] <filepath>\n")
	}
	flag.BoolVar(&clearFlag, "c", false, "clear existing records first")
	flag.BoolVar(&clearFlag, "clear", false, "clear existing records first")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if clearFlag {
		clear()
	}

	load(flag.Args()[0])
}

func clear() {
	fmt.Printf("Clear existing records")
	count, err := owm.ClearCityList()
	if err != nil {
		fmt.Printf("Failed to clear: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Deleted %d records.\n", count)
}

func load(filepath string) {
	fmt.Printf("Load %v\n", filepath)

	count, err := owm.LoadCityList(filepath)
	if err != nil {
		fmt.Printf("Failed to load %v: %s\n", filepath, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Inserted %d records.\n", count)

	count2, err := tritondb.RemoveShiFromJPCities()
	if err != nil {
		fmt.Printf("Failed to update %v: %s\n", filepath, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Updated %d records.\n", count2)
}
