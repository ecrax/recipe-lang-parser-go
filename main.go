package main

import (
	"log"
	"recipe-lang/pkg/generator"
	"recipe-lang/pkg/parser"
)

// TODO: Make this a cli?
// TODO: Generate a PDF with the parsed contents
func main() {
	recipe, err := parser.ParseFromFile("./test.recipe")
	if err != nil {
		log.Fatal(err)
	}
	//
	//marshal, err := json.MarshalIndent(recipe, "", " ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Print(string(marshal))

	out := generator.GenerateHtml(recipe)
	log.Println(string(out))
}
