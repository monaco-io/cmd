package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	input := os.Args[1]
	log.Println(input)
	isTs := len(input) == 10
	if isTs {
		ts, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			log.Fatalln("parse input failed:", err)
			return
		}
		tstr := time.Unix(ts, 0).Format(time.RFC3339)
		fmt.Println("tims is:", tstr)
		return
	}
	// input string
	switch input {
	case "now":
		fmt.Println("unix timestamp now is:", time.Now().Unix())
	default:

	}
}
