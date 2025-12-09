package babel

import (
	"context"
	"encoding/json"
	_ "image/jpeg" // JPEG image support
	_ "image/png"  // PNG image support
	"net/http"
	"strconv"
	customError "swallow-supplier/error"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// TranslateText translate response in target language
func (b *Babel) TranslateText(ctx context.Context, req BabelTextRequest) (res BabelResponse, err error) {
	logger := log.With(b.Service.Logger, "method", "TranslateText")
	level.Info(logger).Log("info", "Service TranslateText")

	response, err := b.Service.Send(
		ctx,
		ServiceName,
		b.Host+"/pre/translate",
		http.MethodPost,
		client.ContentTypeJSON,
		req,
	)

	level.Info(logger).Log("babel translation", "response body", string(response.Body))

	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(b.Ctx, "external_processing_error", "empty response", nil)
	}

	//level.Info(logger).Log("info", "response body ", response.Body)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(b.Ctx, "external_processing_error", "Empty Body", nil)
	}

	res.Code = string(response.Status)
	res.Translations = append(res.Translations, response.Body)
	//level.Info(logger).Log("info", "response body ", res)

	return res, nil
}

// TranslateImage translate source image into target language
func (b *Babel) TranslateImage(ctx context.Context, req BabelImageRequest) (res BabelResponse, err error) {
	logger := log.With(b.Service.Logger, "method", "TranslateImage")
	level.Info(logger).Log("info", "Service TranslateImage")

	response, err := b.Service.Send(
		ctx,
		ServiceName,
		b.Host+"/pre/translate",
		http.MethodPost,
		client.ContentTypeJSON,
		req,
	)

	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(b.Ctx, "external_processing_error", "empty response", nil)
	}

	level.Info(logger).Log("info", "response body ", response.Body)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(b.Ctx, "external_processing_error", "Empty Body", nil)
	}

	return
}

// TranslateImageToText translate source image text into target language
func (b *Babel) TranslateImageToText(ctx context.Context, req BabelImageRequest) (res BabelResponse, err error) {
	logger := log.With(b.Service.Logger, "method", "TranslateImage")
	level.Info(logger).Log("info", "Service TranslateImage")

	response, err := b.Service.Send(
		ctx,
		ServiceName,
		b.Host+"/pre/translate",
		http.MethodPost,
		client.ContentTypeJSON,
		req,
	)

	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(b.Ctx, "external_processing_error", "empty response", nil)
	}

	level.Info(logger).Log("info", "response body ", response.Body)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(b.Ctx, "external_processing_error", "Empty Body", nil)
	}

	return
}

/*func main() {
	// Test downloading image
	url := "https://qa-image6.yanolja.com/leisure/x74fRB2MGr7HGnAr"
	imageBytes, err := DownloadImage(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Downloaded image with %d slices of bytes.\n", len(imageBytes))
}
*/
