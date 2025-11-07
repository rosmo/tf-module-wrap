package main

import (
	"flag"
	"path/filepath"
	"strings"

	tfmodulewrap "github.com/rosmo/tf-module-wrap"
)

func main() {
	var modulePath string
	var moduleVar string
	var ignoreVars string

	flag.StringVar(&modulePath, "module-path", "", "Path containing the module")
	flag.StringVar(&modulePath, "module-var", "", "Variable for the module configuration")
	flag.StringVar(&ignoreVars, "ignore-vars", "", "Variables to ignore")
	flag.Parse()

	ignoreVarsList := make([]string, 0)
	if ignoreVars != "" {
		for _, s := range strings.Split(ignoreVars, ",") {
			ignoreVarsList = append(ignoreVarsList, strings.TrimSpace(s))
		}
	}
	if modulePath == "" {
		flag.PrintDefaults()
		return
	}
	if moduleVar == "" {
		baseName := filepath.Base(modulePath)
		baseName = strings.ReplaceAll(baseName, "-", "_")
		moduleVar = baseName
	}

	err := tfmodulewrap.LoadModule(modulePath, moduleVar, ignoreVarsList)
	if err != nil {
		panic(err)
	}
}
