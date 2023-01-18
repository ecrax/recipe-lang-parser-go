package parser

import (
	"math"
	"os"
	"recipe-lang/internal/duration"
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
	timers := make([]Timer, 0)
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
				timer := Timer{
					Quantity: parseQuantity(groups["timerQuantity"]),
					Name:     groups["timerName"],
					StepType: "timer",
					Units:    parseUnits(groups["timerUnits"]),
				}
				step = append(step, timer)
				timers = append(timers, timer)
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
	}

	// TODO: calculate times
	var totalTime float32
	var cookingTime float32
	var preparationTime float32
	if metadata["total time"] != "" {
		totalTime = duration.ParseDuration(metadata["total time"], duration.Minutes)
	} else {
		if metadata["preparation time"] != "" {
			preparationTime = duration.ParseDuration(metadata["preparation time"], duration.Minutes)
		}
		if metadata["cooking time"] != "" {
			cookingTime = duration.ParseDuration(metadata["cooking time"], duration.Minutes)
		} else {
			i := 0
			// Filter for duplicates
			filteredTimers := util.Filter(timers, func(step Step) bool {
				x := util.FindIndex(len(timers), func(i int) bool {
					t := timers[i]
					return t.Name == step.Name
				}) == i
				i++

				return x
			})

			for _, timer := range filteredTimers {
				cookingTime += duration.ParseDuration(timer.Quantity+timer.Units, duration.Minutes)
			}
		}
		totalTime = preparationTime + cookingTime
	}

	return Recipe{
		title,
		ingredients,
		metadata,
		steps,
		timers,
		Times{
			totalTime,
			float32(math.Round(float64(cookingTime*100)) / 100),
			preparationTime,
		},
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
