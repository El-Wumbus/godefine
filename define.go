package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type Words []struct {
	Word       string      `json:"word"`
	Phonetic   string      `json:"phonetic"`
	Phonetics  []Phonetics `json:"phonetics"`
	Meanings   []Meanings  `json:"meanings"`
	License    License     `json:"license"`
	SourceUrls []string    `json:"sourceUrls"`
}
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Phonetics struct {
	Text      string  `json:"text"`
	Audio     string  `json:"audio"`
	SourceURL string  `json:"sourceUrl"`
	License   License `json:"license"`
}
type Definitions struct {
	Definition string `json:"definition"`
	Synonyms   []any  `json:"synonyms"`
	Antonyms   []any  `json:"antonyms"`
	Example    string `json:"example"`
}
type Meanings struct {
	PartOfSpeech string        `json:"partOfSpeech"`
	Definitions  []Definitions `json:"definitions"`
	Synonyms     []any         `json:"synonyms"`
	Antonyms     []string      `json:"antonyms"`
}

func define(word string) string {
	lword := strings.ToLower(word)

	// Make the API request
	request_url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%v", lword)
	resp, err := http.Get(request_url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "API request failed: %v\n", err)
		os.Exit(74)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API request failed: %v\n", err)
		os.Exit(74)
	}

	// Parse JSON
	var words Words
	err = json.Unmarshal(body, &words)

	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON parse error: %v\n", err)
		os.Exit(65) // :shrug:
	}

	var buffer string
	defined_word := words[0]
	for i := 0; i < len(defined_word.Meanings); i++ {
		meaning := defined_word.Meanings[i]
		buffer += fmt.Sprintf("(%v) %v\n", meaning.PartOfSpeech, meaning.Definitions[0].Definition)
		example := meaning.Definitions[0].Example
		if example != "" {
			buffer += fmt.Sprintf("    Example: '%v'\n", example)
		}
	}
	return fmt.Sprintf("Definitions for '%v':\n%v", word, buffer)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v [word...]\n", path.Base(os.Args[0]))
		os.Exit(64)
	}

	for i := 1; i < len(os.Args); i++ {
		fmt.Println(define(os.Args[i]))
	}
}
