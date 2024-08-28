package utilities

import (
	"errors"
	"fmt"
	"strings"
	"utility/internal/consts"

	"gitlab.com/tuneverse/toolkit/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateLangIdentifier(input string, inputType string, identifier ...string) (string, error) {
	// Check if input or inputType is empty; if so, return an error.
	if utils.IsEmpty(input) || utils.IsEmpty(inputType) {
		return "", errors.New("cannot create language label invalid input or input type")
	}

	// If no identifier is provided, use the default identifier.
	if len(identifier) < 1 {
		identifier = append(identifier, consts.DefaultIdentifier)
	}

	// If the provided identifier is empty, use the default identifier.
	if utils.IsEmpty(identifier[0]) {
		identifier[0] = consts.DefaultIdentifier
	}

	// Define a regular expression pattern to remove non-alphabetical characters.
	// pattern := regexp.MustCompile(`[^a-zA-Z]+`)

	words := strings.Fields(input)
	langIdentifier := strings.ToLower(words[0])

	for i := 1; i < len(words); i++ {
		langIdentifier += cases.Title(language.Und).String(words[i])
	}

	identifier[0] = strings.TrimSpace(identifier[0])

	//langIdentifier := strings.TrimSpace(strings.ReplaceAll(input, " ", ""))

	inputType = strings.TrimSpace(inputType)

	//langIdentifier = strings.ToLower(langIdentifier)
	// Combine the cleaned input, input type, and identifier to generate the final identifier.
	return fmt.Sprintf("%s%s%s", langIdentifier, inputType, identifier[0]), nil
}
