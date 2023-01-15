package parser

type Recipe struct {
	Title       string       `json:"title,omitempty"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
	Metadata    Metadata     `json:"metadata,omitempty"`
	Steps       [][]Step     `json:"steps,omitempty"`
	Times       Times        `json:"times"`
}

type Ingredient = Step

type Metadata = map[string]string

type Step struct {
	Quantity string `json:"quantity,omitempty"`
	Name     string `json:"name,omitempty"`
	StepType string `json:"step_type,omitempty"`
	Units    string `json:"units,omitempty"`
}

type Times struct {
	TotalTime       int `json:"total_time,omitempty"`
	CookingTime     int `json:"cooking_time,omitempty"`
	PreparationTime int `json:"preparation_time,omitempty"`
}
