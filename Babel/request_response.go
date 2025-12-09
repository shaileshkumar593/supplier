package babel

//babel api request for text to text conversion
type BabelTextRequest struct {
	Text       []string `json:"text" validate:"required"`
	TargetLang string   `json:"target_lang" validate:"required"`
}

//babel api request for image to text & image  conversion
type BabelImageRequest struct {
	Images         [][]byte `json:"images"  validate:"required"`
	TargetLanguage string   `json:"target_language" validate:"required"`
}

//babel  api response
type BabelResponse struct {
	Translations []string `json:"translations" validate:"required"`
	Code         string
}
