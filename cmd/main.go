package main

import (
	"fmt"

	"github.com/gwillem/go-simplemon"
)

func main() {
	for name, f := range simplemon.AllChecks {
		fmt.Println("Running check", name)
		if e := f(); e != nil {
			fmt.Println(e)
		}
	}
}
