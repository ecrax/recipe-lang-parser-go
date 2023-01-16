package main

import (
	"encoding/json"
	"fmt"
	"log"
	"recipe-lang/internal/util"
	"recipe-lang/pkg/parser"
)

// TODO: Make this a cli?
// TODO: Generate a PDF with the parsed contents
func main() {
	recipe, err := parser.ParseFromFile("./test.recipe")
	if err != nil {
		log.Fatal(err)
	}

	marshal, err := json.MarshalIndent(recipe, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(marshal))

	//r := util.ParseDuration("2 minutes and 1 hours", "m")
	r := util.ParseDuration("running length: 1hour:20mins", "m")
	log.Println(r)
}
