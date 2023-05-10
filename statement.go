package cgo

import (
	"fmt"
	"strings"
)

type Stmt interface {
	String() string
}

type ReturnStmt struct {
	value interface{}
}

func NewReturnStmt(value interface{}) Stmt {
	return &ReturnStmt{
		value: value,
	}
}

func (s *ReturnStmt) String() string {
	return fmt.Sprintf("return %v", s.value)
}

type FuncCallStmt struct {
	name string
	Args []*Token
}

func NewFuncCallStmt(name string, args []*Token) Stmt {
	return &FuncCallStmt{
		name: name,
		Args: args,
	}
}

func (s *FuncCallStmt) String() string {
	sargs := make([]string, len(s.Args))
	for i, arg := range s.Args {
		sargs[i] = arg.String()
	}
	return fmt.Sprintf("%s(%v)", s.name, strings.Join(sargs, ","))
}

type Func struct {
	Type Type
	Name *Token
	Args []Arg
	Body []Stmt
}

func NewFunc(ttype Type, name *Token, args []Arg, body []Stmt) *Func {
	return &Func{
		Type: ttype,
		Name: name,
		Args: args,
		Body: body,
	}
}
