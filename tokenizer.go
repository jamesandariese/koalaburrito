package koalaburrito

import (
	"fmt"
	"regexp"
)

type Tokenizer struct {
	matchers []*regexp.Regexp
	names    []string
}

// Make a tokenizer
func MakeTokenizer() *Tokenizer {
	return &Tokenizer{}
}

// Add a pattern to the tokenizer
func (t *Tokenizer) AddPattern(pattern, name string) *Tokenizer {
	t.matchers = append(t.matchers, regexp.MustCompile("^"+pattern))
	t.names = append(t.names, name)
	return t
}

// A token with its pattern name and value
type Token struct {
	name  string
	value string
	position int
}

// A tokenization error was found at a specific position.
// This means that the tokenizer wasn't able to match the
// input at t.Position() against any patterns that were
// added to the Tokenizer.
type TokenizationError struct {
	position int
}

// This TokenizationError is an error
// returns true.
func (e *TokenizationError) IsError() bool {
	return true
}

// This Token is not an error
// returns false.
func (t *Token) IsError() bool {
	return false
}

// The position where the tokenization error was encountered.
func (t *TokenizationError) Position() int {
	return t.position
}

// The position where this token starts
func (t *Token) Position() int {
	return t.position
}

// The name of the pattern associated with the token.
func (t *Token) Type() string {
	return t.name
}

// Always return "ERROR" because we couldn't match any pattern
func (t *TokenizationError) Type() string {
	return "ERROR"
}


func (t *Token) String() string {
	return t.value
}

func (t *TokenizationError) String() string {
	return fmt.Sprintf("Error at position %d", t.position)
}

func (t *TokenizationError) Error() string {
	return t.String()
}

type TokenOrError interface {
	IsError() bool
	Position() int
	Type() string
	String() string
}

func (t *Tokenizer) Match(s string, out chan<- TokenOrError) {
	pos := 0
	for s != "" {
		found := false
		for i, m := range t.matchers {
			if loc := m.FindStringIndex(s); loc != nil {
				if loc[0] != 0 {
					panic("A pattern didn't start at 0... inconceivable!")
				}
				out <- &Token{t.names[i], s[:loc[1]], pos+1}
				s = s[loc[1]:]
				pos += loc[1]
				found = true
				break
			}
		}
		if !found {
			out <- &TokenizationError{pos+1}
			break
		}
	}
	close(out)
}

func (t *Tokenizer) MatchAll(s string) ([]*Token, error) {
	c := make(chan TokenOrError)
	ret := []*Token{}
	go t.Match(s, c)
	for t := range c {
		if t.IsError() {
			return []*Token{}, t.(*TokenizationError)
		} else {
			ret = append(ret, t.(*Token))
		}
	}
	return ret, nil
}
