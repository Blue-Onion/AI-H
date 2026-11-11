package main

import (
	"flag"
	"fmt"
)

func main() {
	p := flag.String("p", "", "Value for P")
	m := flag.String("m", "", "Value for m")
	flag.Parse()
	fmt.Println(*p)
	fmt.Println(*m)
}
