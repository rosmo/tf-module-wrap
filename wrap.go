// © Copyright 2025 Taneli Leppä
// SPDX-License-Identifier: Apache-2.0
package tfmodulewrap

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func diagnosticsToError(diags *tfconfig.Diagnostics) error {
	errMsg := ""
	for _, diag := range *diags {
		errMsg += diag.Summary
		if diag.Pos != nil {
			errMsg += fmt.Sprintf(" (%s:%d)", diag.Pos.Filename, diag.Pos.Line)
		}
		if diag.Detail != "" {
			errMsg += "\n" + diag.Detail
		}
		errMsg += "\n"
	}
	return errors.New(errMsg)
}

func printTfValue(v interface{}) string {
	val := ""
	switch dt := v.(type) {
	case map[string]interface{}:
		val += "{\n"
		for kk, vv := range dt {
			val += fmt.Sprintf("  %s = %s\n", kk, printTfValue(vv))
		}
		val += "}"
	case string:
		val += "\"" + strings.ReplaceAll(dt, "\"", "\\\"") + "\""
	case bool:
		if dt {
			val += "true"
		} else {
			val += "false"
		}
	case int:
		val += string(dt)
	case float32:
		val += fmt.Sprintf("%f", dt)
	case float64:
		val += fmt.Sprintf("%f", dt)
	case []interface{}:
		vals := make([]string, 0)
		for vv := range dt {
			vals = append(vals, printTfValue(vv))
		}
		val += fmt.Sprintf("[%s]", strings.Join(vals, ", "))
	case nil:
		val += "null"
	}
	return val
}

func LoadModule(modulePath string, moduleVar string, ignoreVars []string, addDefaultValues bool) error {
	module, _ := tfconfig.LoadModule(modulePath)

	output := ""

	output += fmt.Sprintf("variable \"%s\" {\n", moduleVar)

	// Construct object definition
	output += fmt.Sprintf("  type = object({\n")
	for _, v := range module.Variables {
		if !slices.Contains(ignoreVars, v.Name) {
			varType := v.Type
			if v.Required {
				varType = fmt.Sprintf("optional(%s)", varType)
			}
			if v.Sensitive {
				varType = fmt.Sprintf("sensitive(%s)", varType)
			}
			if v.Description != "" {
				output += fmt.Sprintf("    # %s\n", v.Description)
			}
			output += fmt.Sprintf("    %s = %s", v.Name, v.Type)
			output += fmt.Sprintf("\n")
		}
	}
	output += fmt.Sprintf("  })\n")

	// Add default values
	if addDefaultValues {
		output += "  default = {\n"
		for _, v := range module.Variables {
			if !slices.Contains(ignoreVars, v.Name) {
				if v.Default != nil {
					output += fmt.Sprintf("    %s = %s\n", v.Name, printTfValue(v.Default))

				}
			}
		}
		output += "  }\n"
	}

	if module.Diagnostics.HasErrors() {
		return diagnosticsToError(&module.Diagnostics)
	}
	output += fmt.Sprintf("}\n")

	fmt.Printf("%s\n", hclwrite.Format([]byte(output)))

	fmt.Printf("\n")
	output = fmt.Sprintf("module \"%s\" {\n", moduleVar)
	if !strings.HasPrefix(modulePath, ".") && !strings.HasPrefix(modulePath, "/") {
		output += fmt.Sprintf("  source = \"./%s\"\n\n", modulePath)
	} else {
		output += fmt.Sprintf("  source = \"%s\"\n\n", modulePath)
	}
	for _, v := range module.Variables {
		if !slices.Contains(ignoreVars, v.Name) {
			output += fmt.Sprintf("  %s = var.%s.%s\n", v.Name, moduleVar, v.Name)
		}
	}
	output += "}\n"

	fmt.Printf("%s\n", hclwrite.Format([]byte(output)))

	return nil
}
