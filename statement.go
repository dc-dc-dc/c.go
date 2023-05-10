package cgo

import (
	"fmt"
	"strings"
)

type Stmt interface {
	String() string
}

type ReturnStmt struct {
	Value interface{}
}

func NewReturnStmt(value interface{}) Stmt {
	return &ReturnStmt{
		Value: value,
	}
}

func (s *ReturnStmt) String() string {
	return fmt.Sprintf("return %v", s.Value)
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
	Name string
	Args []Arg
	Body []Stmt
}

func NewFunc(ttype Type, name string, args []Arg, body []Stmt) *Func {
	return &Func{
		Type: ttype,
		Name: name,
		Args: args,
		Body: body,
	}
}
