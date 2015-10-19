package policy

import (
	"reflect"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {

	testCases := []struct {
		input string
		token int
		value yySymType
	}{
		{
			input: "role:member)",
			token: role_check,
			value: yySymType{check: check{match: "member"}},
		},
		{
			input: "rule:blafasel",
			token: rule_check,
			value: yySymType{check: check{match: "blafasel"}},
		},
		{
			input: "user_id:1234",
			token: token_const_check,
			value: yySymType{check: check{key: "user_id", match: "1234"}},
		},
		{
			input: "project_id:%(target.project_id)s",
			token: token_var_check,
			value: yySymType{check: check{key: "project_id", match: "target.project_id"}},
		},
		{
			input: "'a const value':%(target.project_id)s",
			token: const_check,
			value: yySymType{check: check{key: "a const value", match: "target.project_id"}},
		},
		{
			input: "http:https://blafasel.tut",
			token: http_check,
			value: yySymType{check: check{match: "https://blafasel.tut"}},
		},
		{
			input: "and",
			token: and,
		},
		{
			input: "AND",
			token: and,
		},
		{
			input: "(role:tut)",
			token: '(',
		},
		{
			input: ")",
			token: ')',
		},
	}

	for _, c := range testCases {
		reader := strings.NewReader(c.input)
		lexer := NewLexer(reader)
		var s yySymType
		if tt := lexer.Lex(&s); tt != c.token {
			t.Errorf("token type should be %d, got %d", c.token, tt)
		}
		if !reflect.DeepEqual(c.value, s) {
			t.Errorf("Expected value %s, got %s", c.value, s)
		}
	}

}
