package main

import (
	"fmt"
	"log"

	"github.com/CorgiMan/sfmovies/gocode"
)

func main() {
	apidata, err := GetAndParseAPIData()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// 	session.SetMode(mgo.Monotonic, true)

	sfmovies.StoreAPIData(apidata)

	ad2, err := sfmovies.GetLatestAPIData()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println(ad2)

	// err = apidata.WriteToDisc("apidata.json")
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// apidata2, err := sfmovies.ReadFromDisc("apidata.json")
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// fmt.Println(apidata2)
}
