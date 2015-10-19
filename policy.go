package policy

import (
	"errors"
	"fmt"
	"strings"
)

//go:generate go tool yacc -v "" -o parser.go parser.y
//go:generate nex -e -o lexer.go lexer.nex
//go:generate sed -i -e "s/.*NEX_END_OF_LEXER_STRUCT.*/  rules *map[string]Rule/" lexer.go

type rule func(c Context) bool

type Policy interface {
	Enforce(rule string, c Context) bool
}

type policy struct {
	rules map[string]rule
}

func (p *policy) Enforce(rule string, c Context) bool {
	r, ok := p.rules[rule]
	return ok && r(c)

}

func NewPolicy(rules map[string]string) (Policy, error) {
	p := policy{make(map[string]rule, len(rules))}
	for name, str := range rules {
		lexer := NewLexer(strings.NewReader(str))
		lexer.rules = &p.rules
		if yyParse(lexer) != 0 {
			return nil, fmt.Errorf("Failed to parse rule %s: %s", name, lexer.parseResult.(string))
		}
		p.rules[name] = lexer.parseResult.(rule)
	}
	return &p, nil
}

func parseRule(in string) (rule, error) {

	reader := strings.NewReader(in)
	lexer := NewLexer(reader)
	if yyParse(lexer) != 0 {
		return nil, errors.New(lexer.parseResult.(string))

	}
	return lexer.parseResult.(rule), nil
}

type Context struct {
	Token   map[string]string
	Request map[string]string
	Roles   []string
}

func (c Context) hasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (c Context) matchVar(varname, value string) bool {
	v, ok := c.Request[varname]
	return ok && v == value
}

func (c Context) matchToken(key, value string) bool {
	v, ok := c.Token[key]
	return ok && v == value
}

func (c Context) matchTokenAndVar(key, varName string) bool {
	tokenValue, ok := c.Token[key]
	if !ok {
		return false
	}

	if varValue, ok := c.Request[varName]; ok {
		return tokenValue == varValue

	}
	return false
}
