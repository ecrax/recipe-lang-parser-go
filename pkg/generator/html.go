package generator

import (
	"bytes"
	"html/template"
	"recipe-lang/pkg/parser"
)

const tmplPath = "./pkg/generator/templates/"

func GenerateHtml(recipe parser.Recipe) []byte {

	allPaths := []string{tmplPath + "content.tmpl", tmplPath + "header.tmpl", tmplPath + "page.tmpl"}

	templates := template.Must(template.New("").ParseFiles(allPaths...))

	var processed bytes.Buffer
	err := templates.ExecuteTemplate(&processed, "page", recipe.Ingredients)
	if err != nil {
		return nil
	}

	return processed.Bytes()
}
