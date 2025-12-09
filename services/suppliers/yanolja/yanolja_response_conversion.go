package yanolja

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	babel "swallow-supplier/Babel"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func YanoljaResponseConversion(ctx context.Context, response *client.Response, logger log.Logger, err1 error) (res yanolja.Response, err error) {
	res.Code = strconv.Itoa(response.Status)

	if response.Status == 403 {
		res.Code = strconv.Itoa(http.StatusForbidden)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), ServiceName), nil)
		return res, err
	}
	fmt.Println("78")
	if err1 == nil && response.Status == 503 {
		res.Code = strconv.Itoa(http.StatusServiceUnavailable)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-0004", fmt.Sprintf(customError.ErrExternalService.Error(), ServiceName), nil)
		return res, err
	}

	if err1 == nil && response.Status == 504 {
		//fmt.Println("$$$$$$$$$$$$$$ server unavailable $$$$$$$$$$$$$$$$$$$")
		res.Code = strconv.Itoa(http.StatusServiceUnavailable)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-00025", fmt.Sprintf(customError.ErrGatewayTimeout.Error(), ServiceName), nil)
		return res, err
	}

	fmt.Println("79")
	if err1 == nil && response.Body == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Error(logger).Log("err", string(response))
		return res, customError.NewError(ctx, res.Code, "empty response", ServiceName)
	}

	fmt.Println("101")
	if err1 == nil && response.Body != "" && response.Status != http.StatusOK {
		level.Error(logger).Log("error ", response.Body)
		text, _ := ErrorTextConversion(ctx, response, logger)
		res.Body = text
		return res, customError.NewError(ctx, "leisure-api-0012", text, ServiceName)
	}
	fmt.Println("80")
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(ctx, res.Code, "Empty Body", ServiceName)
	}

	fmt.Println("81")
	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	fmt.Println("103")
	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		level.Info(logger).Log(" ..... unmarshaling error .... ")
		res.Code = "500"
		return res, customError.NewError(ctx, res.Code, "unmarshaling error ", nil)
	}

	fmt.Println("104")

	level.Info(logger).Log(" yanolja response for waitfororder ")

	if res.Code != "200" {
		jsonstr, err2 := json.Marshal(res.Body)
		if err2 != nil {
			level.Error(logger).Log("error", "Error in marshaling data:", err2)
			res.Body = "Error in marshaling data"
			res.Code = "500"
			return res, err2
		}
		var resBody yanolja.ResponseBody
		// Unmarshal the JSON data into the struct
		err2 = json.Unmarshal(jsonstr, &resBody)
		if err2 != nil {
			level.Error(logger).Log("error", "Error in unmarshaling data:", err2)

			res.Body = "Error in unmarshaling data :"
			res.Code = "500"
			return res, err2
		}
		return res, customError.NewErrorCustom(ctx, res.Code, resBody.Detail, resBody.Message, response.Status, ServiceName)
	}

	return res, nil
}

func ErrorTextConversion(ctx context.Context, response *client.Response, logger log.Logger) (text string, err error) {

	// Declare a variable to hold the map
	result := make(map[string]interface{})

	// Prepare the array for the translation request
	arry := make([]string, 0)

	// Unmarshal the JSON string into the map
	err = json.Unmarshal([]byte(response.Body), &result)
	if err != nil {
	}
	// Check if "body" exists in the map
	body, bodyExists := result["body"].(map[string]interface{})
	if !bodyExists {
		level.Error(logger).Log("error ", "body is empty")
		return "", fmt.Errorf("body is empty")
	}

	// Access the "message" key in the "body" map
	message, messageExists := body["message"].(string)
	if !messageExists {
		level.Error(logger).Log("error ", "message is empty")
		return "", fmt.Errorf("message is empty")
	}
	arry = append(arry, message)

	// Access the "detail" key in the "body" map
	detail, detailExists := body["detail"].(string)
	if !detailExists {
		level.Error(logger).Log("error ", "detail is empty")
		return "", fmt.Errorf("detail is empty")
	}
	arry = append(arry, detail)
	// Access the "properties" key in the "body" map
	_, propertiesExists := body["properties"].(map[string]interface{})
	if !propertiesExists {
		level.Error(logger).Log("error ", "properties  is empty")
		return "", fmt.Errorf("properties is empty")
	}

	level.Info(logger).Log("Info", "call to babel")
	// Print the accessed values
	output := babel.BabelTextRequest{
		Text:       arry,
		TargetLang: "EN-US",
	}

	level.Info(logger).Log("info", "request to babel api", output)

	// Call external translation service
	babelsvc, _ := babel.New(ctx)
	productTranslation, err := babelsvc.TranslateText(ctx, output)
	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raised an error", err)
		return "", fmt.Errorf("babel api error: %w", err)
	}
	level.Info(logger).Log("Translated Text ", productTranslation.Translations)
	Text, err := MergeText(productTranslation.Translations, logger)
	return Text, err
}

func MergeText(data []string, logger log.Logger) (mergedText string, err error) {
	// Iterate over the slice of JSON strings
	var mergedMessage string
	for _, val := range data {
		fmt.Println("Processing JSON:", val)

		// Declare a variable to hold the result as a map
		var result map[string]interface{}

		// Unmarshal the JSON into the map
		err = json.Unmarshal([]byte(val), &result)
		if err != nil {
			level.Error(logger).Log("error", fmt.Sprintf("Error unmarshalling JSON: %v", err))
			return "", err
		}

		// Extract translations field
		translations, ok := result["translations"].([]interface{})
		if !ok {
			level.Error(logger).Log("Error: ", "translations field is not a slice")
			return "", fmt.Errorf("'translations' field is not a slice")
		}

		// Iterate over the translations and extract "text" field
		for _, translation := range translations {
			// Assert that the item in translations is a map
			translationMap, ok := translation.(map[string]interface{})
			if !ok {
				level.Error(logger).Log("Error: ", " translation is not a map")
				return "", fmt.Errorf("translation is not a map")
			}

			// Extract the "text" field
			text, ok := translationMap["text"].(string)
			if !ok {
				level.Error(logger).Log("Error: ", "text field is missing or not a string")
				return "", fmt.Errorf("'text' field is missing or not a string")
			}

			// Append the text to the merged message
			mergedMessage += text + " "
		}
	}

	// Trim any extra spaces at the end
	mergedMessage = strings.TrimSpace(mergedMessage)

	fmt.Println("->->->->->->->->  Final merged message:", mergedMessage)

	// Return the merged text
	return mergedMessage, nil
}
