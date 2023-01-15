package main

import (
	"encoding/json"
	"fmt"
	"log"
	"recipe-lang/pkg/parser"
)

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
}
