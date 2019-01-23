package config

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

// Parse parses the input and returs the result.
func Parse(input []byte) (map[string]interface{}, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input        []byte
	pos          int
	result       map[string]interface{}
	err          error
	stringCloser byte
}

func newLex(input []byte) *lex {
	return &lex{
		input: input,
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	return l.scanNormal(lval)
}

func (l *lex) scanNormal(lval *yySymType) int {
	for b := l.next(); b != 0; b = l.next() {
		switch {
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"' || b == '\'':
			l.stringCloser = b
			return l.scanString(lval)
		case unicode.IsDigit(rune(b)) || b == '+' || b == '-':
			l.backup()
			return l.scanNum(lval)
		case unicode.IsLetter(rune(b)):
			l.backup()
			return l.scanLiteral(lval)
		default:
			return int(b)
		}
	}
	return 0
}

var escape = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func (l *lex) scanString(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := l.next(); b != 0; b = l.next() {
		switch {
		case b == '\\':
			b2 := escape[l.next()]
			if b2 == 0 {
				return LexError
			}
			buf.WriteByte(b2)
		case b == l.stringCloser:
			lval.val = buf.String()
			return String
		default:
			buf.WriteByte(b)
		}
	}
	return LexError
}

func (l *lex) scanNum(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case unicode.IsDigit(rune(b)):
			buf.WriteByte(b)
		case strings.IndexByte(".+-eE", b) != -1:
			buf.WriteByte(b)
		default:
			l.backup()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return LexError
			}
			lval.val = val
			return Number
		}
	}
}

func (l *lex) scanLiteral(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case unicode.IsLetter(rune(b)) || b == '_':
			buf.WriteByte(b)
		default:
			l.backup()
			lval.val = buf.String()
			return Literal
		}
	}
}

func (l *lex) backup() {
	if l.pos == -1 {
		return
	}
	l.pos--
}

func (l *lex) next() byte {
	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
		return 0
	}
	l.pos++
	return l.input[l.pos-1]
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}
