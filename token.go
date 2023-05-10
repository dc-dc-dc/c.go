package cgo

import "fmt"

type TokenType int
type Type int

const (
	TokenOperand TokenType = iota
	TokenCperand
	TokenOcurly
	TokenCcurly
	TokenComma
	TokenSemiColon
	TokenName
	TokenNumber
	TokenString
	TokenReturn
)

const (
	TypeInvalid Type = iota
	TypeInt
	TypeVoid
	TypeChar
)

var (
	literalTokens = map[rune]TokenType{
		'(': TokenOperand,
		')': TokenCperand,
		'{': TokenOcurly,
		'}': TokenCcurly,
		',': TokenComma,
		';': TokenSemiColon,
	}
)

func (t Type) String() string {
	switch t {
	case TypeChar:
		return "char"
	case TypeInt:
		return "int"
	case TypeVoid:
		return "void"
	default:
		return "unknown"
	}
}

func (t TokenType) String() string {
	switch t {
	case TokenOperand:
		return "("
	case TokenCperand:
		return ")"
	case TokenOcurly:
		return "{"
	case TokenCcurly:
		return "}"
	case TokenComma:
		return ","
	case TokenSemiColon:
		return ";"
	case TokenString:
		return "string"
	case TokenNumber:
		return "number"
	case TokenName:
		return "name"
	case TokenReturn:
		return "return"
	default:
		return "unknown"
	}
}

type Location struct {
	filePath string
	row      int
	col      int
}

func NewLocation(filePath string, row, col int) *Location {
	return &Location{
		filePath: filePath,
		row:      row,
		col:      col,
	}
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.filePath, l.row+1, l.col+1)
}

type Token struct {
	TokenType TokenType
	Value     interface{}
	location  *Location
}

func NewToken(tokenType TokenType, value interface{}, location *Location) *Token {
	return &Token{
		TokenType: tokenType,
		Value:     value,
		location:  location,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s: %v", t.TokenType.String(), t.Value)
}
