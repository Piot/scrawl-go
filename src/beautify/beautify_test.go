package beautify

import (
	"fmt"
	"strings"
	"testing"

	"github.com/piot/scrawl-go/src/runestream"
	"github.com/piot/scrawl-go/src/token"
	"github.com/piot/scrawl-go/src/tokenize"
)

func setupTokenizer(text string) *tokenize.Tokenizer {
	ioReader := strings.NewReader(text)
	runeReader := runestream.NewRuneReader(ioReader)
	tokenizer := tokenize.NewTokenizer(runeReader)
	return tokenizer
}

func setupBeautify(x string) (string, error) {
	t := setupTokenizer(x)
	var tokens []token.Token
	for {
		tok, tokErr := t.ReadNext()
		fmt.Printf("token:%v\n", tok)
		if tokErr != nil {
			return "", tokErr
		}
		if tok == nil {
			break
		}
		tokens = append(tokens, tok)
	}
	builder := &strings.Builder{}
	Write(builder, tokens)
	output := builder.String()
	fmt.Printf("output %v\n", output)
	return output, nil
}

func check(t *testing.T, test string, expected string) {
	beautified, err := setupBeautify(test)
	if err != nil {
		t.Fatal(err)
	}
	if expected != beautified {
		t.Errorf("mismatch. Expected %q but got %q", expected, beautified)
	}
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
