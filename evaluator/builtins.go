package evaluator

import (
	"github.com/ruegerj/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"puts": object.GetBuiltinByName("puts"),
	"len":  object.GetBuiltinByName("len"),
	// Array builtins
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),
}
