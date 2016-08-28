package koalaburrito

import (
	"testing"
)

func TestTokenization(t *testing.T) {
        tokenizer := MakeTokenizer()
        tokenizer.AddPattern(`([-+]?([0-9]+\.[0-9]*|[0-9]*\.[0-9]+|[0-9]+))`, "NUMBER")
        tokenizer.AddPattern(`\+`, "ADD")
        tokenizer.AddPattern(`-`, "SUBTRACT")
        tokenizer.AddPattern(`/`, "DIVIDE")
        tokenizer.AddPattern(`\*`, "MULTIPLY")
        tokenizer.AddPattern(`\s+`, "WHITESPACE")
        c := make(chan TokenOrError)
        go tokenizer.Match("5.5 2 1 1 +1 -1.1 +/-*", c)
	runOne := func(tokenType, tokenString string) {
		tok := <-c
		if tok.Type() != tokenType {
			t.Errorf("Token type did not match: expected %v but got %#v", tokenType, tok)
		}
		if tok.String() != tokenString {
			t.Errorf("Token string did not match: expected %v but got %#v", tokenString, tok)
		}
	}
	runOne("NUMBER", "5.5")
	runOne("WHITESPACE", " ")
	runOne("NUMBER", "2")
	runOne("WHITESPACE", " ")
	runOne("NUMBER", "1")
	runOne("WHITESPACE", " ")
	runOne("NUMBER", "1")
	runOne("WHITESPACE", " ")
	runOne("NUMBER", "+1")
	runOne("WHITESPACE", " ")
	runOne("NUMBER", "-1.1")
	runOne("WHITESPACE", " ")
	runOne("ADD", "+")
	runOne("DIVIDE", "/")
	runOne("SUBTRACT", "-")
	runOne("MULTIPLY", "*")
}

func TestTokenizationError(t *testing.T) {
        tokenizer := MakeTokenizer()
        tokenizer.AddPattern(`\s+`, "WHITESPACE")
        c := make(chan TokenOrError)
        go tokenizer.Match("    5", c)
	runOne := func(tokenType, tokenString string, position int) {
		tok := <-c
		if tok.Type() != tokenType {
			t.Errorf("Token type did not match: expected %v but got %#v", tokenType, tok)
		}
		if tok.String() != tokenString {
			t.Errorf("Token string did not match: expected %v but got %#v", tokenString, tok)
		}
		if tok.Position() != position {
			t.Errorf("Token position was expected to be %v but instead it was %v", position, tok.Position())
		}
	}
	runOne("WHITESPACE", "    ", 1)
	runOne("ERROR", "Error at position 5", 5)
}
