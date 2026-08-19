package translator

import (
	"fmt"
	"strconv"
)

type TokenKind uint8

const (
	TokenWord TokenKind = iota
	TokenNumber
	TokenString
)

type Token struct {
	Kind  TokenKind
	Text  string
	Value int32
}

func tokenize(source string) ([]Token, error) {
	tokens := make([]Token, 0)
	i := 0

	for i < len(source) {
		if isWhitespace(source[i]) {
			i++
			continue
		}

		if isComment(source[i]) {
			i = readUntil(i, '\n', source)
			continue
		}

		if isStringStart(source, i) {
			i += 2
			if i < len(source) && isWhitespace(source[i]) {
				i++
			}
			start := i
			i = readUntil(i, '"', source)
			if i >= len(source) {
				return nil, fmt.Errorf("unterminated string literal")
			}
			tokens = append(tokens, Token{
				Kind: TokenString,
				Text: source[start:i],
			})
			i++
			continue
		}

		start := i
		i = readWordEnd(start, source)
		lexeme := source[start:i]

		if value, err := strconv.ParseInt(lexeme, 10, 32); err == nil {
			tokens = append(tokens, Token{
				Kind:  TokenNumber,
				Text:  lexeme,
				Value: int32(value),
			})
		} else {
			tokens = append(tokens, Token{
				Kind: TokenWord,
				Text: lexeme,
			})
		}

	}

	return tokens, nil
}

func readUntil(i int, until byte, source string) int {
	for i < len(source) && source[i] != until {
		i++
	}
	return i
}

func readWordEnd(i int, source string) int {
	for i < len(source) && !isWhitespace(source[i]) {
		i++
	}
	return i
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isComment(b byte) bool {
	return b == '\\'
}

func isStringStart(source string, i int) bool {
	return i+1 < len(source) && source[i] == 's' && source[i+1] == '"'
}
