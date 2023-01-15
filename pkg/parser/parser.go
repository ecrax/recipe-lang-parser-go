package parser

import (
	"os"
	"recipe-lang/internal/util"
	"regexp"
	"strings"
)

var commentRegex = regexp.MustCompile(`(?m)--.*`)
var blockCommentRegex = regexp.MustCompile(`(?m)\s*\[-.*?-]\s*`)

const titleRegex = `(?m)^#\s*\s*(?P<title>.+)`
const metadataRegex = `(?m)^>>\s*(?P<key>.+?):\s*(?P<value>.+)`
const multiwordIngredientRegex = `(?m)#(?P<mIngredientName>[^@#~[]+?){(?P<mIngredientQuantity>.*?)(?:_(?P<mIngredientUnits>[^}]+?))?}`
const singlewordIngredientRegex = `(?m)#(?P<sIngredientName>[^\s]+)`
const timerIngredientRegex = `(?m)@(?P<timerName>.*?)(?:{(?P<timerQuantity>.*?)(?:_(?P<timerUnits>.+?))?})`

var allRegex = [...]string{titleRegex, metadataRegex, multiwordIngredientRegex, singlewordIngredientRegex, timerIngredientRegex}
var tokens = strings.Join(allRegex[:], "|")
var tokensRegex = regexp.MustCompile(tokens)

func Parse(source string) (Recipe, error) {
	ingredients := make([]Ingredient, 0)
	steps := make([][]Step, 0)
	metadata := Metadata{}
	title := ""

	// Remove all comments
	source = commentRegex.ReplaceAllString(source, " ")
	source = blockCommentRegex.ReplaceAllString(source, " ")

	lines := strings.Split(source, "\n")

	// Filter empty lines
	lines = util.Filter(lines, func(s string) bool {
		return len(strings.Trim(s, " \n\r")) > 0
	})

	pos := 0
	for _, line := range lines {
		step := make([]Step, 0)
		matches := util.FindNamedMatches(tokensRegex, line)

		for _, match := range matches {
			//log.Printf("%v | %s", match.Index, match.Context)
			groups := match.Group

			if groups["title"] != "" {
				title = groups["title"]
			}

			if match.Index != 0 && match.Index < pos {
				step = append(step, Step{StepType: "text", Name: line[0:match.Index]})
			}

			if pos < match.Index {
				step = append(step, Step{StepType: "text", Name: line[pos:match.Index]})
			}

			if groups["key"] != "" && groups["value"] != "" {
				metadata[strings.ToLower(strings.Trim(groups["key"], " \r\n"))] = strings.Trim(groups["value"], " \r\n")
			}

			if groups["sIngredientName"] != "" {
				ingredient := Ingredient{
					Quantity: "",
					Name:     groups["sIngredientName"],
					StepType: "ingredient",
					Units:    "",
				}
				ingredients = append(ingredients, ingredient)
				step = append(step, ingredient)
			}

			if groups["mIngredientName"] != "" {
				ingredient := Ingredient{
					Quantity: parseQuantity(groups["mIngredientQuantity"]),
					Name:     groups["mIngredientName"],
					StepType: "ingredient",
					Units:    parseUnits(groups["mIngredientUnits"]),
				}
				ingredients = append(ingredients, ingredient)
				step = append(step, ingredient)
			}

			if groups["timerQuantity"] != "" {
				step = append(step, Step{
					Quantity: parseQuantity(groups["timerQuantity"]),
					Name:     groups["timerName"],
					StepType: "timer",
					Units:    parseUnits(groups["timerUnits"]),
				})
			}

			pos = match.Index + len(match.Context)
		}

		if pos < len(line) {
			step = append(step, Step{
				Quantity: "",
				Name:     line[pos:],
				StepType: "text",
				Units:    "",
			})
		}

		if len(step) > 0 {
			steps = append(steps, step)
		}
		//log.Println()
	}

	// TODO: calculate times

	return Recipe{
		Title:       title,
		Ingredients: ingredients,
		Metadata:    metadata,
		Steps:       steps,
		Times:       Times{},
	}, nil
}

func parseQuantity(quantity string) string {
	return strings.Trim(quantity, " \r\n")
}

func parseUnits(units string) string {
	return strings.Trim(units, " \r\n")
}

func ParseFromFile(path string) (Recipe, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Recipe{}, err
	}

	return Parse(string(data))
}
