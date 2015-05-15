package main

import (
	"fmt"

	"github.com/CorgiMan/sfmovies"
)

func main() {
	apidata, err := sfmovies.GetAndParseTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = apidata.WriteToDisc("apidata.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	apidata2, err := sfmovies.ReadFromDisc("apidata.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(apidata2)
}
