package monkey

import (
	"fmt"
	"testing"
)

func init() {

}

func getTokens(ch chan item) []item {
	tokens := []item{}
	for i := range ch {
		tokens = append(tokens, i)
	}
	return tokens
}

func printTokens(tokens []item) {
	for _, i := range tokens {
		fmt.Println(i.String())
	}
}

func TestLex1(t *testing.T) {
	content := `let a = 10;`
	l := NewLexer(content)
	go l.run()
	tokens := getTokens(l.Items)
	printTokens(tokens)
	t.Log("--all--")
	return
}

func TestFloat1(t *testing.T) {
	content := `let a = 0.9;`
	l := NewLexer(content)
	go l.run()
	tokens := getTokens(l.Items)
	printTokens(tokens)
	t.Log("--all--")
	return
}

func TestOperator(t *testing.T) {
	content := `let a = 100 + 0.9 - 1.44;`
	l := NewLexer(content)
	go l.run()
	tokens := getTokens(l.Items)
	printTokens(tokens)
	t.Log("--all--")
	return

}
