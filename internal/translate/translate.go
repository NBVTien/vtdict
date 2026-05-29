package translate

import (
	"fmt"

	"resty.dev/v3"
)

type responseData struct {
	TranslatedText string `json:"translatedText"`
}

type response struct {
	ResponseData responseData `json:"responseData"`
}

// Uses MyMemory free translation API — no key needed, 1000 req/day free
func Translate(text, targetLang string) (string, error) {
	client := resty.New()
	var result response

	langPair := fmt.Sprintf("en|%s", targetLang)
	_, err := client.R().
		SetResult(&result).
		SetQueryParam("q", text).
		SetQueryParam("langpair", langPair).
		Get("https://api.mymemory.translated.net/get")

	if err != nil {
		return "", err
	}
	if result.ResponseData.TranslatedText == "" {
		return "", fmt.Errorf("translation failed")
	}
	return result.ResponseData.TranslatedText, nil
}
