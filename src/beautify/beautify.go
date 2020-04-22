package beautify

import (
	"fmt"
	"io"
	"strings"

	"github.com/piot/scrawl-go/src/token"
)

type OutputStream struct {
	writer io.Writer
	line   int
	col    int
	indent int
}

func (o *OutputStream) writeSymbol(symbol token.SymbolToken) {
	fmt.Fprintf(o.writer, "%v", symbol.Symbol)
	o.col++
}

func (o *OutputStream) writeOperator(symbol token.OperatorToken) {
	fmt.Fprintf(o.writer, "%v", symbol.Operator)
	o.col++
}

func (o *OutputStream) writeString(stringToken token.StringToken) {
	fmt.Fprintf(o.writer, "'%v'", stringToken.Text())
	o.col++
}

func (o *OutputStream) writeComment(commentToken token.CommentToken) {
	fmt.Fprintf(o.writer, "# %s", commentToken.Text())
	o.col++
}

func (o *OutputStream) writeNumber(numberToken token.NumberToken) {
	fmt.Fprintf(o.writer, "%d", numberToken.Integer())
	o.col++
}

func (o *OutputStream) writeIndent() {
	for i := 0; i < o.indent; i++ {
		fmt.Fprint(o.writer, "  ")
	}
	//o.col++
}

func (o *OutputStream) writeStartScope() {
	o.writeLineDelimiter()
	o.indent++
}

func (o *OutputStream) writeStartMeta() {
	fmt.Fprint(o.writer, "[")
}

func (o *OutputStream) writeEndMeta() {
	fmt.Fprint(o.writer, "]")
}

func (o *OutputStream) writeLineDelimiter() {
	fmt.Fprintf(o.writer, "\n")
	o.col = 0
	o.line++
}
func (o *OutputStream) writeEndScope() {
	o.indent--
}

func writeToken(o *OutputStream, tok token.Token) {
	_, wasStartScope := tok.(token.StartScopeToken)
	if wasStartScope {
		o.writeStartScope()
		return
	}
	_, wasEndScope := tok.(token.EndScopeToken)
	if wasEndScope {
		o.writeEndScope()
		return
	}
	_, wasLineDelimiter := tok.(token.LineDelimiterToken)
	if wasLineDelimiter {
		o.writeLineDelimiter()
		return
	}
	if o.col != 0 {
		fmt.Fprint(o.writer, " ")
	} else {
		o.writeIndent()
	}
	_, wasStartMeta := tok.(token.StartMetaDataToken)
	if wasStartMeta {
		o.writeStartMeta()
		return
	}
	_, wasEndMeta := tok.(token.EndMetaDataToken)
	if wasEndMeta {
		o.writeEndMeta()
		return
	}
	symbolToken, wasSymbol := tok.(token.SymbolToken)
	if wasSymbol {
		o.writeSymbol(symbolToken)
		return
	}
	stringToken, wasString := tok.(token.StringToken)
	if wasString {
		o.writeString(stringToken)
		return
	}
	commentToken, wasComment := tok.(token.CommentToken)
	if wasComment {
		o.writeComment(commentToken)
		return
	}
	numberToken, wasNumber := tok.(token.NumberToken)
	if wasNumber {
		o.writeNumber(numberToken)
		return
	}

	operatorToken, wasOperatorToken := tok.(token.OperatorToken)
	if wasOperatorToken {
		o.writeOperator(operatorToken)
		return
	}

	panicString := fmt.Sprintf("unknown token %v %T", tok, tok)
	panic(panicString)
}

type BeautifyFlags uint8

const (
	Normal BeautifyFlags = 1 << iota
	DiscardComments
)

func Write(writer io.Writer, tokens []token.Token, flags BeautifyFlags) {
	o := &OutputStream{writer: writer}
	for _, tok := range tokens {
		if flags&DiscardComments != 0 {
			_, wasComment := tok.(token.CommentToken)
			if wasComment {
				continue
			}
		}

		writeToken(o, tok)
	}
}

func Process(tokens []token.Token, flags BeautifyFlags) string {
	builder := &strings.Builder{}
	Write(builder, tokens, flags)
	output := builder.String()
	return output
}
