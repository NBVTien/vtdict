package dictionary

import (
	"fmt"

	"resty.dev/v3"
)

type Phonetic struct {
	Text string `json:"text"`
}

type Definition struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type Result struct {
	Word      string     `json:"word"`
	Phonetics []Phonetic `json:"phonetics"`
	Meanings  []Meaning  `json:"meanings"`
}

func Lookup(word string) ([]Result, error) {
	client := resty.New()
	var results []Result
	var apiErr map[string]interface{}

	resp, err := client.R().
		SetResult(&results).
		SetError(&apiErr).
		Get(fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("word not found")
	}
	return results, nil
}

func FirstDefinition(results []Result) string {
	for _, r := range results {
		for _, m := range r.Meanings {
			if len(m.Definitions) > 0 {
				return m.Definitions[0].Definition
			}
		}
	}
	return ""
}
