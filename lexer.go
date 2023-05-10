package cgo

import (
	"fmt"
	"strconv"
	"unicode"
)

type Lexer struct {
	filePath string
	source   string
	cur      int
	bol      int
	row      int
}

func NewLexer(filePath, source string) *Lexer {
	return &Lexer{
		filePath: filePath,
		source:   source,
		cur:      0,
		bol:      0,
		row:      0,
	}
}

func (l *Lexer) IsEmpty() bool {
	return l.cur >= len(l.source)
}

func (l *Lexer) ChopChar() {
	if !l.IsEmpty() {
		x := l.source[l.cur]
		l.cur += 1
		if x == '\n' {
			l.bol = l.cur
			l.row += 1
		}
	}
}

func (l *Lexer) TrimLeft() {
	for !l.IsEmpty() && unicode.IsSpace(rune(l.source[l.cur])) {
		l.ChopChar()
	}
}

func (l *Lexer) DropLine() {
	for !l.IsEmpty() && l.source[l.cur] != '\n' {
		l.ChopChar()
	}
	if !l.IsEmpty() {
		l.ChopChar()
	}
}

func (l *Lexer) NextToken() *Token {
	l.TrimLeft()
	for !l.IsEmpty() && l.source[l.cur] == '#' {
		l.DropLine()
		l.TrimLeft()
	}
	if l.IsEmpty() {
		return nil
	}
	loc := NewLocation(l.filePath, l.row, l.cur-l.bol)
	first := l.source[l.cur]
	if unicode.IsLetter(rune(first)) {
		start := l.cur
		for !l.IsEmpty() && (unicode.IsLetter(rune(l.source[l.cur])) || unicode.IsDigit(rune(l.source[l.cur]))) {
			l.ChopChar()
		}
		val := l.source[start:l.cur]
		return NewToken(TokenName, val, loc)
	}

	if tokenType, ok := literalTokens[rune(first)]; ok {
		l.ChopChar()
		return NewToken(tokenType, first, loc)
	}

	if first == '"' {
		l.ChopChar()
		start := l.cur
		for !l.IsEmpty() && l.source[l.cur] != '"' {
			l.ChopChar()
		}
		if !l.IsEmpty() {
			val := l.source[start:l.cur]
			l.ChopChar()
			return NewToken(TokenString, val, loc)
		}
		fmt.Printf("error: unclosed string literal at %s", loc.String())
		return nil
	}

	if unicode.IsDigit(rune(first)) {
		start := l.cur
		for !l.IsEmpty() && unicode.IsDigit(rune(l.source[l.cur])) {
			l.ChopChar()
		}
		raw := l.source[start:l.cur]
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			fmt.Printf("error: failed to parse int at loc: %s, err: %s", loc.String(), err.Error())
		}
		return NewToken(TokenNumber, val, loc)
	}
	return nil
}
