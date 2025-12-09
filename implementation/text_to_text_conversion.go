package implementation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	babel "swallow-supplier/Babel"
	customError "swallow-supplier/error"
	model "swallow-supplier/mongo/domain/yanolja"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/net/context"
)

func TextToTextConversion(ctx context.Context, logger log.Logger, product model.Product) (productrec model.Product, err error) {
	jsonData, err := json.Marshal(product)
	if err != nil {
		level.Error(logger).Log("error", "marshalling error", err)
		return productrec, fmt.Errorf("error marshaling product: %w", err)
	}

	// Prepare the array for the translation request
	arry := []string{string(jsonData)}
	output := babel.BabelTextRequest{
		Text:       arry,
		TargetLang: "EN-US",
	}

	//level.Info(logger).Log("info", "request to babel api", output)

	// Call external translation service
	babelsvc, _ := babel.New(ctx)
	productTranslation, err := babelsvc.TranslateText(ctx, output)
	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raised an error", err)
		return productrec, fmt.Errorf("yanolja client error: %w", err)
	}

	// Unmarshal the translated JSON back into a map
	translatedProduct := make(map[string][]map[string]interface{})
	err = json.Unmarshal([]byte(productTranslation.Translations[0]), &translatedProduct)
	if err != nil {
		level.Error(logger).Log("error", "response unmarshal error", err)
		return productrec, fmt.Errorf("error unmarshaling translated product: %w", err)
	}

	//level.Info(logger).Log("info", "translated response ", translatedProduct)

	// Extract the "text" field from the translations map
	ll, ok := translatedProduct["translations"][0]["text"]
	if !ok {
		level.Info(logger).Log("Info", "error extracting translated text", nil)
		return productrec, fmt.Errorf("translated text not found")
	}

	// Check if `ll` is a string containing JSON and unmarshal again if necessary
	var jsonString string
	if str, ok := ll.(string); ok {
		jsonString = str
	} else {
		level.Error(logger).Log("error", "Translated text is not in the expected string format", nil)
		return productrec, fmt.Errorf("translated text not a string")
	}

	//  sanitize invalid JSON escapes (\R → \\R)
	jsonString = sanitizeInvalidEscapes(jsonString)

	//  (optional) compact it to remove any stray whitespace / trailing commas
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(jsonString)); err == nil {
		jsonString = buf.String()
	} else {
		level.Warn(logger).Log("msg", "could not compact JSON", "err", err)
		// we continue with the un-compacted version
	}

	// Unmarshal the JSON string into the Product struct
	var translatedProductStruct model.Product
	err = json.Unmarshal([]byte(jsonString), &translatedProductStruct)
	if err != nil {
		level.Error(logger).Log("error", "error unmarshaling translated text to Product struct: ", err)
		return productrec, fmt.Errorf("error unmarshaling translated text to struct: %w", err)
	}

	if product.ProductID != translatedProductStruct.ProductID {
		level.Info(logger).Log("Info", "Error unmarshaling translated text to Product struct: %v", err)
		return productrec, customError.NewError(ctx, "leisure-api-0005", fmt.Sprintf("requested data is not in proper json format, %v", err), "UpsertProduct")
	}

	return translatedProductStruct, nil
}

// precompile once at package level
var invalidEscapePattern = regexp.MustCompile(`\\([^"\\/bfnrtu])`)

// sanitizeInvalidEscapes will turn `\R` into `\\R`, etc.
func sanitizeInvalidEscapes(s string) string {
	return invalidEscapePattern.ReplaceAllString(s, `\\$1`)
}

/*
func TextToTextConversion(ctx context.Context, logger log.Logger, product model.Product) (model.Product, error) {

	// 1) Marshal original product to JSON

	jsonData, err := json.Marshal(product)

	if err != nil {

		level.Error(logger).Log("error", "marshalling product", "err", err)

		return model.Product{}, fmt.Errorf("failed to marshal input product: %w", err)

	}

	// 2) Call the Babel translation API

	req := babel.BabelTextRequest{

		Text: []string{string(jsonData)},

		TargetLang: "EN-US",
	}

	babelSvc, err := babel.New(ctx)

	if err != nil {

		level.Error(logger).Log("error", "initializing babel client", "err", err)

		return model.Product{}, fmt.Errorf("failed to initialize translation client: %w", err)

	}

	resp, err := babelSvc.TranslateText(ctx, req)

	if err != nil {

		level.Error(logger).Log("error", "translation API call failed", "err", err)

		return model.Product{}, fmt.Errorf("translation service error: %w", err)

	}

	if len(resp.Translations) == 0 {

		level.Error(logger).Log("error", "no translations returned from API", nil)

		return model.Product{}, fmt.Errorf("translation service returned no data")

	}

	// 3) Extract the raw translation

	raw := resp.Translations[0]

	// 4) If the service wrapped the JSON in quotes, unquote it

	//    e.g. "\"{ ... }\""

	if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {

		unq, uqErr := strconv.Unquote(raw)

		if uqErr == nil {

			raw = unq

		} else {

			// not critical—just log and continue with raw

			level.Warn(logger).Log("msg", "failed to unquote translation", "err", uqErr)

		}

	}

	// 5) Sanitize *all* invalid JSON escapes in one shot

	//    turns \R → \\R, \X → \\X for any X not in the JSON escape set

	const invalidEscapeRe = `\\([^"\\/bfnrtu])`

	re := regexp.MustCompile(invalidEscapeRe)

	sanitized := re.ReplaceAllString(raw, `\\$1`)

	// 6) (Optional) Compact it to remove stray whitespace/trailing commas

	var buf bytes.Buffer

	if compErr := json.Compact(&buf, []byte(sanitized)); compErr == nil {

		sanitized = buf.String()

	} else {

		level.Warn(logger).Log("msg", "could not compact JSON", "err", compErr)

	}

	// 7) Now unmarshal straight into your Product struct

	var translated model.Product

	if err := json.Unmarshal([]byte(sanitized), &translated); err != nil {

		level.Error(logger).Log("error", "unmarshaling translated product", "err", err)

		return model.Product{}, fmt.Errorf("failed to unmarshal translated product: %w", err)

	}

	// 8) Verify that critical fields match

	if translated.ProductID != product.ProductID {

		level.Error(logger).Log("error", "productID mismatch after translation", "got", translated.ProductID, "want", product.ProductID)

		return model.Product{}, customError.NewError(ctx,

			"leisure-api-0005",

			"translated productID does not match original",

			"TextToTextConversion",
		)

	}

	return translated, nil

}

var (
	// precompile regex to catch any invalid JSON escape like \R, \X, etc.
	invalidEscapePattern = regexp.MustCompile(`\\([^"\\/bfnrtu])`)
)

// TextToTextConversion calls an external translation API to convert
// all fields of product into English, then unmarshals back into model.Product.
// It logs and returns detailed errors on unknown fields, type mismatches, or syntax errors.
func TextToTextConversion(
	ctx context.Context,
	logger log.Logger,
	product model.Product,
) (model.Product, error) {
	// 1) Marshal input product to JSON
	inputJSON, err := json.Marshal(product)
	if err != nil {
		level.Error(logger).Log("msg", "failed to marshal input product", "err", err)
		return model.Product{}, fmt.Errorf("marshal input product: %w", err)
	}

	// 2) Prepare and send translation request
	req := babel.BabelTextRequest{
		Text:       []string{string(inputJSON)},
		TargetLang: "EN-US",
	}
	babelSvc, err := babel.New(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize babel client", "err", err)
		return model.Product{}, fmt.Errorf("init translation client: %w", err)
	}

	resp, err := babelSvc.TranslateText(ctx, req)
	if err != nil {
		level.Error(logger).Log("msg", "translation API call failed", "err", err)
		return model.Product{}, fmt.Errorf("translation service error: %w", err)
	}
	if len(resp.Translations) == 0 {
		level.Error(logger).Log("msg", "no translations returned")
		return model.Product{}, fmt.Errorf("translation service returned no data")
	}

	raw := resp.Translations[0]

	// 3) Unquote if wrapped in Go‐style quotes
	if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {
		if unq, uqErr := strconv.Unquote(raw); uqErr == nil {
			raw = unq
		} else {
			level.Warn(logger).Log("msg", "failed to unquote translation", "err", uqErr)
		}
	}

	// 4) Sanitize invalid JSON escapes (\R → \\R, etc.)
	sanitized := invalidEscapePattern.ReplaceAllString(raw, `\\$1`)

	// 5) Compact JSON to remove stray whitespace/trailing commas
	var buf bytes.Buffer
	if compErr := json.Compact(&buf, []byte(sanitized)); compErr == nil {
		sanitized = buf.String()
	} else {
		level.Warn(logger).Log("msg", "could not compact JSON", "err", compErr)
	}

	// 6) Decode into Product with strict error reporting
	var translated model.Product
	dec := json.NewDecoder(strings.NewReader(sanitized))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&translated); err != nil {
		// unknown field
		if strings.HasPrefix(err.Error(), "json: unknown field ") {
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			level.Error(logger).Log("msg", "unknown JSON field", "field", field)
			return model.Product{}, fmt.Errorf("unknown JSON field: %s", field)
		}
		// type mismatch
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			level.Error(logger).Log(
				"msg", "type mismatch in JSON",
				"struct", typeErr.Struct,
				"field", typeErr.Field,
				"expected", typeErr.Type,
				"value", typeErr.Value,
				"offset", typeErr.Offset,
			)
			return model.Product{}, fmt.Errorf(
				"type mismatch for %s.%s: expected %v but got %s at offset %d",
				typeErr.Struct, typeErr.Field, typeErr.Type, typeErr.Value, typeErr.Offset,
			)
		}
		// syntax error
		var synErr *json.SyntaxError
		if errors.As(err, &synErr) {
			level.Error(logger).Log(
				"msg", "JSON syntax error",
				"offset", synErr.Offset,
				"err", synErr.Error(),
			)
			return model.Product{}, fmt.Errorf("malformed JSON at byte %d: %v", synErr.Offset, synErr.Error())
		}
		// other errors
		level.Error(logger).Log("msg", "failed to decode translated product", "err", err)
		return model.Product{}, fmt.Errorf("decode translated product: %w", err)
	}

	// 7) Ensure critical field consistency
	if translated.ProductID != product.ProductID {
		level.Error(logger).Log(
			"msg", "product ID mismatch after translation",
			"got", translated.ProductID,
			"expected", product.ProductID,
		)
		return model.Product{}, customError.NewError(
			ctx,
			"leisure-api-0005",
			fmt.Sprintf("translated productID %q does not match original %q", translated.ProductID, product.ProductID),
			"TextToTextConversion",
		)
	}

	return translated, nil
}
*/
