package beautify

import (
	"fmt"
	"strings"
	"testing"

	"github.com/piot/scrawl-go/src/parser"
	"github.com/piot/scrawl-go/src/token"
)

func fetchAllTokens(x string) ([]token.Token, error) {
	t := parser.SetupTokenizer(x)
	tokens, readErr := t.ReadAll()

	return tokens, readErr
}

func setupBeautify(x string) ([]token.Token, string, error) {
	tokens, err := fetchAllTokens(x)
	if err != nil {
		return nil, "", err
	}
	builder := &strings.Builder{}
	Write(builder, tokens)
	output := builder.String()
	fmt.Printf("output %v\n", output)
	return tokens, output, nil
}

func check(t *testing.T, test string, expected string) {
	_, beautified, err := setupBeautify(test)
	if err != nil {
		t.Fatal(err)
	}
	if expected != beautified {
		t.Errorf("mismatch. Expected %q but got %q", expected, beautified)
	}
}

func compareTokens(actualTokens []token.Token, expectedTokens []token.Token) error {
	if len(actualTokens) != len(expectedTokens) {
		return fmt.Errorf("not even the same length %v vs %v", len(actualTokens), len(expectedTokens))
	}

	for index, expectedToken := range expectedTokens {
		actualToken := actualTokens[index]
		if !expectedToken.IsEqual(actualToken) {
			return fmt.Errorf("token mismatch %v and %v", expectedToken, actualToken)
		}
	}

	return nil
}

func compareTokensHelper(t *testing.T, actualTokens []token.Token, expectedTokens []token.Token) {
	err := compareTokens(actualTokens, expectedTokens)
	if err != nil {
		t.Error(err)
	}
}

func checkReverse(t *testing.T, test string, expected string) {
	tokens, beautified, err := setupBeautify(test)
	if err != nil {
		t.Fatal(err)
	}
	if expected != beautified {
		t.Errorf("mismatch. Expected %q but got %q", expected, beautified)
	}
	beautifiedTokens, fetchErr := fetchAllTokens(beautified)
	if fetchErr != nil {
		t.Error(fetchErr)
	}

	compareTokensHelper(t, beautifiedTokens, tokens)
}

func TestAnything(t *testing.T) {
	check(t, "    hello    Goodbye", "hello Goodbye")
	check(t, "    hello    Goodbye    -2313", "hello Goodbye -2313")
	check(t, "    hello    Goodbye    -2313 	'hi'", "hello Goodbye -2313 'hi'")
	check(t,
		`
Something 3212
  other 03
`,
		`Something 3212
  other 3
`)
}

func TestAnythingReverse(t *testing.T) {
	checkReverse(t, "    hello    Goodbye", "hello Goodbye")
	checkReverse(t, "    hello    Goodbye    -2313", "hello Goodbye -2313")
	checkReverse(t, "    hello    Goodbye    -2313 	'hi'", "hello Goodbye -2313 'hi'")
	checkReverse(t,
		`
Something 3212
  other 03
`,
		`Something 3212
  other 3
`)
}
