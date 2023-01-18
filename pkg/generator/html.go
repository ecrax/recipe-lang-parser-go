package generator

import (
	"bytes"
	"html/template"
	"log"
	"recipe-lang/pkg/parser"
)

const tmplPath = "./pkg/generator/templates/"

func GenerateHtml(recipe parser.Recipe) []byte {

	allPaths := []string{tmplPath + "content.html", tmplPath + "header.html", tmplPath + "page.html"}

	templates := template.Must(template.New("").ParseFiles(allPaths...))

	var processed bytes.Buffer
	err := templates.ExecuteTemplate(&processed, "page", recipe)
	if err != nil {
		log.Fatal(err)
	}

	return processed.Bytes()
}
