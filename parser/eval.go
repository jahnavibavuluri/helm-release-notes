package parser

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func BuildEvalContext(vars map[string]TFVariable, locals map[string]cty.Value) *hcl.EvalContext {
	varMap := map[string]cty.Value{}
	for name, variable := range vars {
		val := strings.Trim(variable.Default, `"`)
		varMap[name] = cty.StringVal(val)
	}

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var":   cty.ObjectVal(varMap),
			"local": cty.ObjectVal(locals),
		},
	}
}
