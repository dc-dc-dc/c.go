package cgo

import (
	"fmt"
)

func ExpectToken(lexer *Lexer, tokenType ...TokenType) *Token {
	token := lexer.NextToken()
	if token == nil {
		fmt.Printf("error: expected token %s\n", tokenType)
		return nil
	}
	for _, ttype := range tokenType {
		if token.TokenType == ttype {
			return token
		}
	}
	fmt.Printf("error: expected token %s, but got %s\n", tokenType, token.String())
	return nil
}

func ParseType(lexer *Lexer) Type {
	return ParseTypeToken(ExpectToken(lexer, TokenName))
}

func ParseTypeToken(token *Token) Type {
	if token == nil {
		return TypeInvalid
	}
	switch token.Value {
	case "int":
		return TypeInt
	case "void":
		return TypeVoid
	case "char":
		return TypeChar
	default:
		fmt.Printf("error: unknown type at %s, got: %c\n", token.location.String(), token.Value)
		return TypeInvalid
	}
}

type Arg struct {
	Name string
	Type Type
}

func ParseFuncArgList(lexer *Lexer) []Arg {
	if ExpectToken(lexer, TokenOperand) == nil {
		return nil
	}
	args := []Arg{}
	expr := ExpectToken(lexer, TokenName, TokenCperand)
	if expr == nil {
		return nil
	}
	if expr.TokenType == TokenCperand {
		return args
	}
	ttype := ParseTypeToken(expr)
	if ttype == TypeInvalid {
		return nil
	}
	name := ExpectToken(lexer, TokenName)
	if name == nil {
		return nil
	}
	args = append(args, Arg{name.Value.(string), ttype})
	for {
		expr = ExpectToken(lexer, TokenComma, TokenCperand, TokenName)
		// fmt.Printf("expr: %v\n", expr)
		if expr == nil {
			return nil
		}
		if expr.TokenType == TokenComma {
			continue
		}
		if expr.TokenType == TokenCperand {
			break
		}
		ttype := ParseTypeToken(expr)
		if ttype == TypeInvalid {
			return nil
		}
		name := ExpectToken(lexer, TokenName)
		if name == nil {
			return nil
		}
		args = append(args, Arg{name.Value.(string), ttype})
	}
	return args
}

func ParseArgList(lexer *Lexer) []*Token {
	if ExpectToken(lexer, TokenOperand) == nil {
		return nil
	}
	args := []*Token{}
	expr := ExpectToken(lexer, TokenString, TokenNumber, TokenCperand)
	if expr == nil {
		return nil
	}
	if expr.TokenType == TokenCperand {
		return args
	}

	args = append(args, expr)
	for {
		expr = ExpectToken(lexer, TokenComma, TokenCperand)
		if expr == nil {
			return nil
		}
		if expr.TokenType == TokenCperand {
			break
		}
		expr = ExpectToken(lexer, TokenString, TokenNumber)
		if expr == nil {
			return nil
		}
		args = append(args, expr)
	}

	return args
}

func ParseBlock(lexer *Lexer) []Stmt {
	if ExpectToken(lexer, TokenOcurly) == nil {
		return nil
	}

	block := []Stmt{}

	for {
		name := ExpectToken(lexer, TokenName, TokenCcurly)
		if name == nil {
			return nil
		}
		if name.TokenType == TokenCcurly {
			break
		}
		if name.Value == "return" {
			expr := ExpectToken(lexer, TokenNumber, TokenString)
			if expr == nil {
				return nil
			}
			block = append(block, NewReturnStmt(expr.Value))
		} else {
			arglist := ParseArgList(lexer)
			if arglist == nil {
				return nil
			}
			block = append(block, NewFuncCallStmt(name.Value.(string), arglist))
		}
		if ExpectToken(lexer, TokenSemiColon) == nil {
			return nil
		}
	}

	return block
}

func ParseFunction(lexer *Lexer) *Func {
	if lexer.IsEmpty() {
		return nil
	}

	// return type of the function
	ttype := ParseType(lexer)
	if ttype == TypeInvalid {
		return nil
	}
	// name of the function
	name := ExpectToken(lexer, TokenName)
	if name == nil {
		return nil
	}

	// Get the args list from the function def
	args := ParseFuncArgList(lexer)

	// The body of statements
	body := ParseBlock(lexer)

	return NewFunc(ttype, name.Value.(string), args, body)
}
