package main

import (
	"log"
	"os"
	"recipe-lang/pkg/generator"
	"recipe-lang/pkg/parser"
)

// TODO: Make this a cli?
// TODO: Generate a PDF with the parsed contents
// https://github.com/SebastiaanKlippert/go-wkhtmltopdf
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
	err = os.WriteFile("recipe.html", out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
