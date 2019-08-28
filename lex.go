package monkey

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type itemType int

type item struct {
	typ itemType
	val string
}

/*
let x = 55;
let y = 5 + 5;
*/
const (
	itemErr itemType = iota
	itemLet
	itemNum
	itemFloat
	itemVar
	itemAssignment
	itemOperator
	itemEnter
	itemEOL // EOL = ";"
	itemEOF
)

/*
const (
	OperatorTBL = []string{
		"!",
		">",
		"<",
		"==",
		">=",
		"<=",
		"!=",
	}
)
*/

func (i item) String() string {
	switch i.typ {
	case itemLet:
		return "{token:itemLet}"
	case itemNum:
		return fmt.Sprintf("{token:itemNum, val:'%s'}", i.val)
	case itemFloat:
		return fmt.Sprintf("{token:itemFloat, val:'%s'}", i.val)
	case itemVar:
		return fmt.Sprintf("{token:itemVar, val:'%s'}", i.val)
	case itemAssignment:
		return "{token:itemAssignment}"
	case itemOperator:
		return fmt.Sprintf("{token:itemOperator, val:'%s'}", i.val)
	case itemEnter:
		return "{token:itemEnter}"
	case itemEOL:
		return "{token:itemEOL}"
	case itemEOF:
		return "{token:itemEOF}"
	default:
		return fmt.Sprintf("{Token:%d, val:'%s'}", i.typ, i.val)
	}
	return fmt.Sprintf("{unknowToken:%d, val:'%s'}", i.typ, i.val)
}

type lexer struct {
	content string
	start   int
	pos     int
	width   int
	Items   chan item
}

func NewLexer(content string) *lexer {
	return &lexer{
		content: content,
		Items:   make(chan item),
	}
}

func (l *lexer) next() rune {
	// EOF
	if int(l.pos) >= len(l.content) {
		l.width = 0
		return -1
	}
	r, w := utf8.DecodeRuneInString(l.content[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) rollBack() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.rollBack()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// push the token
func (l *lexer) emit(typ itemType) {
	l.Items <- item{typ, l.content[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.Items)
}

/* help functions */
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isEnter(r rune) bool {
	return r == '\n'
}

func isNumberic(r rune) bool {
	return unicode.IsDigit(r)
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isOperatorChar(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '=' || r == '>' || r == '<' || r == '!'
}

/* state functions  */
type stateFn func(*lexer) stateFn

func lexText(l *lexer) stateFn {
	for {
		r := l.peek()
		switch {
		case r == -1:
			l.emit(itemEOF)
			return nil
		case isLetter(r):
			l.next()
			return lexWord
		case isNumberic(r):
			l.next()
			if r == '0' {
				return lexFloat
			} else {
				return lexNum
			}
		case isSpace(r):
			l.next()
			l.ignore()
		case r == ';':
			l.next()
			l.emit(itemEOL)
			return lexText
		case isEnter(r):
			l.next()
			l.emit(itemEnter)
			return lexText
		case isOperatorChar(r):
			return lexOperator
		default:
			l.next()
			l.emit(itemErr)
			return nil
		}

	}
}

func lexWord(l *lexer) stateFn {
	for {
		r := l.peek()
		switch {
		case isLetter(r) || isNumberic(r) || r == '_':
			l.next()
		default:
			switch l.content[l.start:l.pos] {
			case "let":
				l.emit(itemLet)
			default:
				l.emit(itemVar)
			}
			return lexText
		}
	}
}

func lexNum(l *lexer) stateFn {
	for {
		r := l.peek()
		switch {
		case isNumberic(r):
			l.next()
		case r == '.':
			return lexFloat
		case isLetter(r):
			l.emit(itemErr)
			return nil
		default:
			l.emit(itemNum)
			return lexText
		}
	}
}

func lexFloat(l *lexer) stateFn {
	l.next() // read dot aka '.'
	for {
		r := l.peek()
		switch {
		case isNumberic(r):
			l.next()
		case isLetter(r):
			l.emit(itemErr)
			return nil
		default:
			l.emit(itemFloat)
			return lexText
		}
	}
}

func lexOperator(l *lexer) stateFn {
	head := true
	for {
		r := l.peek()
		switch {
		case r == '=' && head:
			l.next()
			if nr := l.peek(); !isOperatorChar(nr) {
				l.emit(itemAssignment)
				return lexText
			}
		case isOperatorChar(r):
			l.next()
		default:
			l.emit(itemOperator)
			return lexText
		}
	}
}
