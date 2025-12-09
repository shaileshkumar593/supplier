package yanolja

type ImageUrl struct {
	ListOfUrl      []string `json:"listOfUrl"  validate:"required"`
	TargetLanguage string   `json:"targetLanguage" validate:"required"`
}
